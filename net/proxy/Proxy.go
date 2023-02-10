package proxy

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"time"
)

type Protocol string

const (
	HTTP   = "http"
	SOCKS4 = "socks4"
	SOCKS5 = "socks5"
)

type (
	Proxy struct {
		Protocol Protocol
		Host     string
		Auth     *Auth
		Timeout  time.Duration
	}
	Auth struct {
		Username string
		Password string
	}
)

func (p Proxy) GetString() (key string) {
	if p.Auth != nil {
		key = string(p.Protocol) + "://" + p.Auth.Username + ":" + p.Auth.Password + "@" + p.Host
	} else {
		key = string(p.Protocol) + "://" + p.Host
	}
	return
}

func (p Proxy) Dial() func(string) (net.Conn, error) {
	switch p.Protocol {
	case HTTP:
		return func(s string) (net.Conn, error) {
			return p.dialHTTP(s)
		}
	case SOCKS4:
		return func(s string) (net.Conn, error) {
			return p.dialSOCKS4(s)
		}
	case SOCKS5:
		return func(s string) (net.Conn, error) {
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
	case "http":
		proxy.Protocol = HTTP
	case "socks4":
		proxy.Protocol = SOCKS4
	case "socks5":
		proxy.Protocol = SOCKS5
	default:
		return nil, fmt.Errorf("unknown proxy protocol %s", uri.Scheme)
	}

	proxy.Host = uri.Host
	usr := uri.User.Username()
	passwd, _ := uri.User.Password()

	if usr != "" || passwd != "" {
		if len(usr) > 255 || len(passwd) > 255 {
			return nil, errors.New("invalid user name or password")
		}
		proxy.Auth = &Auth{
			Username: usr,
			Password: passwd,
		}
	}

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
