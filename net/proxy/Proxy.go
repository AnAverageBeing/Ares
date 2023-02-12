package proxy

import (
	"fmt"
	"net"
	"net/url"
	"time"
)

const (
	SOCKS4 = "socks4"
	SOCKS5 = "socks5"
)

type Proxy struct {
	Protocol Protocol
	Host     string
	Auth     *url.Userinfo
	Timeout  time.Duration
}

type Protocol string

func (p Proxy) GetString() string {
	if p.Auth != nil {
		return string(p.Protocol) + "://" + p.Auth.String() + "@" + p.Host
	}
	return string(p.Protocol) + "://" + p.Host
}

func (p Proxy) Dial() func(string) (conn net.Conn, err error) {
	switch p.Protocol {
	case SOCKS4:
		return func(s string) (conn net.Conn, err error) {
			return p.dialSOCKS4(s)
		}
	case SOCKS5:
		return func(s string) (conn net.Conn, err error) {
			return p.dialSOCKS5(s)
		}
	}
	return nil
}

func New(proxyUri string) (*Proxy, error) {
	uri, err := url.Parse(proxyUri)
	if err != nil {
		return nil, err
	}

	proxy := &Proxy{}

	switch uri.Scheme {
	case SOCKS4:
		proxy.Protocol = SOCKS4
	case SOCKS5:
		proxy.Protocol = SOCKS5
	default:
		return nil, fmt.Errorf("unknown proxy protocol %s", uri.Scheme)
	}

	proxy.Host = uri.Host
	proxy.Auth = uri.User

	query := uri.Query()

	timeout := query.Get("timeout")
	if timeout != "" {
		var err error
		proxy.Timeout, err = time.ParseDuration(timeout)
		if err != nil {
			return nil, err
		}
	}

	return proxy, nil
}
