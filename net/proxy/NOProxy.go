package proxy

import "net"

func (p *Proxy) dialNOProxy(target string) (net.Conn, error) {
	return net.DialTimeout("tpc", target, p.Timeout)
}
