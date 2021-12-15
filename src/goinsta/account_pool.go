package goinsta

import (
	"container/list"
	"makemoney/common/log"
	"sync"
	"time"
)

type AccountPoolt struct {
	accounts     []*Instagram
	notAvailable *list.List
	Available    *list.List
	Cooling      *list.List
	avalLock     sync.Mutex
	noAvalLock   sync.Mutex
	coolingLock  sync.Mutex
	checkTimer   *time.Ticker
}

var AccountPool AccountPoolt

var CallBackCheckAccount func(inst *Instagram) bool

func InitAccountPool(accounts []*Instagram) {
	AccountPool.Cooling = list.New()
	AccountPool.Available = list.New()
	AccountPool.notAvailable = list.New()

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

	if CallBackCheckAccount != nil {
		AccountPool.checkTimer = time.NewTicker(time.Second * 10)
		go CheckAccount()
	}
}

func CheckAccount() {
	for _ = range AccountPool.checkTimer.C {
		AccountPool.coolingLock.Lock()
		for item := AccountPool.Cooling.Front(); item != nil; item = item.Next() {
			inst := item.Value.(*Instagram)
			if CallBackCheckAccount(inst) {
				AccountPool.ReleaseOne(inst)
				next := item.Next()
				AccountPool.Cooling.Remove(item)
				item = next
				log.Info("check account: %s cooling finish", inst.User)
			}
		}
		AccountPool.coolingLock.Unlock()
	}
}

func (this *AccountPoolt) GetOneBlock() *Instagram {
	for true {
		if this.Available.Len() > 0 {
			inst := this.GetOneNoWait()
			if inst != nil {
				return inst
			}
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (this *AccountPoolt) GetOneNoWait() *Instagram {
	this.avalLock.Lock()
	defer this.avalLock.Unlock()

	ret := this.Available.Front()
	if ret == nil {
		return nil
	}
	this.Available.Remove(ret)
	return ret.Value.(*Instagram)
}

func (this *AccountPoolt) GetOne(block bool) *Instagram {
	if block {
		return this.GetOneBlock()
	} else {
		return this.GetOneNoWait()
	}
}

func (this *AccountPoolt) ReleaseOne(insta *Instagram) {
	this.avalLock.Lock()
	defer this.avalLock.Unlock()

	this.Available.PushBack(insta)
}

func (this *AccountPoolt) CoolingOne(inst *Instagram) {
	this.coolingLock.Lock()
	defer this.coolingLock.Unlock()

	this.Cooling.PushBack(inst)
}

func (this *AccountPoolt) BlackOne(insta *Instagram) {
	this.noAvalLock.Lock()
	defer this.noAvalLock.Unlock()

	this.notAvailable.PushBack(insta)
	SaveInstToDB(insta)
}
