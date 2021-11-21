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
	AccountPool.accounts = accounts
	for idx := range accounts {
		AccountPool.Available.PushBack(accounts[idx])
	}
}

func (this *AccountPoolt) GetOne() *Instagram {
	this._lock.Lock()
	defer this._lock.Unlock()

	ret := this.Available.Front()
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
