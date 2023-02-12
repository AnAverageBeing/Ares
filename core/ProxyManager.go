package core

import "Ares/net/proxy"

type ProxyManager struct {
	proxies []*proxy.Proxy
	atIndex int
}

func (manager *ProxyManager) Add(proxy *proxy.Proxy) {
	manager.proxies = append(manager.proxies, proxy)
}

func (manager *ProxyManager) Remove(proxy *proxy.Proxy) {
	for i, p := range manager.proxies {
		if p.GetString() == proxy.GetString() {
			manager.proxies = append(manager.proxies[:i], manager.proxies[i+1:]...)
			return
		}
	}
}

func (manager *ProxyManager) Length() (length int) {
	length = len(manager.proxies)
	return
}

func (manager *ProxyManager) GetNext() *proxy.Proxy {
	manager.atIndex = manager.atIndex + 1
	if manager.atIndex >= len(manager.proxies) {
		manager.atIndex = 0
	}
	return manager.proxies[manager.atIndex]
}

