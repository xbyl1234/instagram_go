package goinsta

import (
	"container/list"
	"sync"
)

type AccountPoolt struct {
	accounts     []*Instagram
	notAvailable list.List
	Available    list.List
	_lock        sync.Mutex
}

var AccountPool AccountPoolt

func InitAccountPool(accounts []*Instagram) {

	for idx := range accounts {
		if accounts[idx].IsLogin && accounts[idx].Status == "" {
			AccountPool.Available.PushBack(accounts[idx])
		} else {
			AccountPool.notAvailable.PushBack(accounts[idx])
		}
	}
	AccountPool.accounts = make([]*Instagram, AccountPool.Available.Len())
	var index = 0
	for item := AccountPool.Available.Front(); item != nil; item = item.Next() {
		AccountPool.accounts[index] = item.Value.(*Instagram)
		index++
	}
}

func (this *AccountPoolt) GetOne() *Instagram {
	this._lock.Lock()
	defer this._lock.Unlock()

	ret := this.Available.Front()
	if ret == nil {
		return nil
	}
	this.Available.Remove(ret)
	return ret.Value.(*Instagram)
}

func (this *AccountPoolt) ReleaseOne(insta *Instagram) {
	this._lock.Lock()
	defer this._lock.Unlock()

	this.Available.PushBack(insta)
}

func (this *AccountPoolt) BlackOne(insta *Instagram) {
	this._lock.Lock()
	defer this._lock.Unlock()

	this.notAvailable.PushBack(insta)
}
