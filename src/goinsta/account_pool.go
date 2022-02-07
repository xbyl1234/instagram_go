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
	Cur          *list.Element
	avalLock     sync.Mutex
	noAvalLock   sync.Mutex
	coolingLock  sync.Mutex
	checkTimer   *time.Ticker
}

var AccountPool *AccountPoolt

var CallBackCheckAccount func(inst *Instagram) bool = nil

func InitAccountPool(accounts []*Instagram) {
	AccountPool = &AccountPoolt{}
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

	log.Info("init account pool available count: %d, bad account :%d",
		AccountPool.Available.Len(), AccountPool.notAvailable.Len())
}

func CheckAccount() {
	for range AccountPool.checkTimer.C {
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

func (this *AccountPoolt) GetOneBlock(OperName string) *Instagram {
	for true {
		if this.Available.Len() > 0 {
			log.Warn("try require %s account", OperName)
			inst := this.GetOneNoWait(OperName)
			if inst != nil {
				return inst
			}
		}
		time.Sleep(time.Second * 10)
	}
	return nil
}

func (this *AccountPoolt) GetOneNoWait(OperName string) *Instagram {
	this.avalLock.Lock()
	defer this.avalLock.Unlock()
	if this.Cur == nil {
		this.Cur = this.Available.Front()
		if this.Cur == nil {
			return nil
		}
	}
	var oldCurl = this.Cur
	for true {
		inst := this.Cur.Value.(*Instagram)
		lastCur := this.Cur
		this.Cur = this.Cur.Next()

		if !inst.IsSpeedLimit(OperName) {
			this.Available.Remove(lastCur)
			return inst
		}
		if this.Cur == nil {
			this.Cur = this.Available.Front()
			if this.Cur == nil {
				return nil
			}
		}
		if this.Cur == oldCurl {
			break
		}
	}
	return nil
}

func (this *AccountPoolt) ReleaseOne(insta *Instagram) {
	this.avalLock.Lock()
	defer this.avalLock.Unlock()
	if insta.IsBad() {
		this.notAvailable.PushBack(insta)
	} else {
		this.Available.PushBack(insta)
	}
	SaveInstToDB(insta)
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
