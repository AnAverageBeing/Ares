package proxy

import (
	"errors"
	"net"
)

func (p *Proxy) dialSOCKS5(target string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", p.Host, p.Timeout)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	var req requestBuilder

	version := byte(5) // socks version 5
	method := byte(0)  // method 0: no authentication (only anonymous access supported for now)
	if p.Auth != nil {
		method = 2 // method 2: username/password
	}

	// version identifier/method selection request
	req.add(
		version, // socks version
		1,       // number of methods
		method,
	)

	resp, err := p.sendReceive(conn, req.Bytes())
	if err != nil {
		return nil, err
	} else if len(resp) != 2 {
		return nil, errors.New("server does not respond properly")
	} else if resp[0] != 5 {
		return nil, errors.New("server does not support Socks 5")
	} else if resp[1] != method {
		return nil, errors.New("socks method negotiation failed")
	}
	if p.Auth != nil {
		version := byte(1) // user/password version 1
		req.Reset()
		req.add(
			version,                    // user/password version
			byte(len(p.Auth.Username)), // length of username
		)
		req.add([]byte(p.Auth.Username)...)
		req.add(byte(len(p.Auth.Password)))
		req.add([]byte(p.Auth.Password)...)
		resp, err := p.sendReceive(conn, req.Bytes())
		if err != nil {
			return nil, err
		} else if len(resp) != 2 {
			return nil, errors.New("server does not respond properly")
		} else if resp[0] != version {
			return nil, errors.New("server does not support user/password version 1")
		} else if resp[1] != 0 { // not success
			return nil, errors.New("user/password login failed")
		}
	}

	// detail request
	host, port, err := splitHostPort(target)
	if err != nil {
		return nil, err
	}
	req.Reset()
	req.add(
		5,               // version number
		1,               // connect command
		0,               // reserved, must be zero
		3,               // address type, 3 means domain name
		byte(len(host)), // address length
	)
	req.add([]byte(host)...)
	req.add(
		byte(port>>8), // higher byte of destination port
		byte(port),    // lower byte of destination port (big endian)
	)
	resp, err = p.sendReceive(conn, req.Bytes())
	if err != nil {
		return conn, err
	} else if len(resp) != 10 {
		return nil, errors.New("proxy does not respond properly")
	} else if resp[1] != 0 {
		return nil, errors.New("can't complete SOCKS5 connection")
	}
	return conn, nil
}
