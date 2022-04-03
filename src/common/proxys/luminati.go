package proxys

//
//type LuminatiPool struct {
//	ProxyImpl
//	allCount  int
//	allProxys map[string]*Proxy
//	ProxyList []*Proxy
//	proxyLock sync.Mutex
//	path      string
//	dumpsPath string
//}
//
//func InitLuminatiPool(path string) (ProxyImpl, error) {
//	var pool = &LuminatiPool{}
//	pool.path = path
//	pool.dumpsPath = strings.ReplaceAll(path, ".json", "_dumps.json")
//	if common.PathExists(pool.dumpsPath) {
//		path = pool.dumpsPath
//	}
//
//	var ProxyMap map[string]*Proxy
//	err := common.LoadJsonFile(path, &ProxyMap)
//	if err != nil {
//		return nil, err
//	}
//
//	ProxyList := make([]*Proxy, len(ProxyMap))
//	var index = 0
//	for _, vul := range ProxyMap {
//		if vul.BlackType != BlacktypeNoblack {
//			continue
//		}
//		ProxyList[index] = vul
//		index++
//	}
//	if len(ProxyList) == 0 {
//		return nil, &common.MakeMoneyError{ErrStr: "no proxy", ErrType: common.PorxyError}
//	}
//
//	pool.ProxyList = ProxyList[:index]
//	pool.allCount = len(ProxyMap)
//	pool.allProxys = ProxyMap
//	return pool, nil
//}
//
//func (this *LuminatiPool) GetNoRisk(busy bool, used bool) *Proxy {
//	this.proxyLock.Lock()
//	defer this.proxyLock.Unlock()
//	for true {
//		index := math_rand.Intn(len(this.ProxyList))
//		if this.ProxyList[index].BlackType == BlacktypeNoblack {
//			if busy {
//				if this.ProxyList[index].IsBusy {
//					continue
//				}
//				this.ProxyList[index].IsBusy = true
//			}
//
//			if used {
//				if this.ProxyList[index].IsUsed {
//					continue
//				}
//				this.ProxyList[index].IsUsed = true
//			}
//
//			return this.ProxyList[index]
//		}
//	}
//	return nil
//}
//
//func (this *LuminatiPool) Get(id string) *Proxy {
//	this.proxyLock.Lock()
//	defer this.proxyLock.Unlock()
//	return this.allProxys[id]
//}
//
//func (this *LuminatiPool) Black(proxy *Proxy, _type BlackType) {
//	this.proxyLock.Lock()
//	defer this.proxyLock.Unlock()
//	proxy.BlackType = _type
//	this.remove(proxy)
//	this.Dumps()
//}
//
//func (this *LuminatiPool) Remove(proxy *Proxy) {
//	this.proxyLock.Lock()
//	defer this.proxyLock.Unlock()
//	this.remove(proxy)
//}
//
//func (this *LuminatiPool) remove(proxy *Proxy) {
//	find := false
//	var index int
//	for index = range this.ProxyList {
//		if this.ProxyList[index] == proxy {
//			find = true
//			break
//		}
//	}
//	if find {
//		delete(this.allProxys, proxy.ID)
//		this.ProxyList = append(this.ProxyList[:index], this.ProxyList[index+1:]...)
//	}
//}
//
//func (this *LuminatiPool) Dumps() {
//	err := common.Dumps(this.dumpsPath, this.allProxys)
//	if err != nil {
//		log.Error("dumps proxy pool error:%v", err)
//	}
//}
