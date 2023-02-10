package proxy

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"time"
)

type requestBuilder struct {
	bytes.Buffer
}

func (b *requestBuilder) add(data ...byte) {
	_, _ = b.Write(data)
}

func (c *Proxy) sendReceive(conn net.Conn, req []byte) (resp []byte, err error) {
	if c.Timeout > 0 {
		conn.SetWriteDeadline(time.Now().Add(c.Timeout))
	}
	_, err = conn.Write(req)
	if err != nil {
		return
	}
	resp, err = c.readAll(conn)
	return
}

func (c *Proxy) readAll(conn net.Conn) (resp []byte, err error) {
	if c.Timeout > 0 {
		conn.SetReadDeadline(time.Now().Add(c.Timeout))
	}
	var n int
	resp = make([]byte, 1024)
	n, err = conn.Read(resp)
	if err != nil {
		return nil, err
	}
	resp = resp[:n]
	return
}

func lookupIPv4(host string) (net.IP, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		if ip4 := ip.To4(); ip4 != nil {
			return ip4, nil
		}
	}
	return nil, fmt.Errorf("no IPv4 address found for host: %s", host)
}

func splitHostPort(addr string) (host string, port uint16, err error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, err
	}
	portInt, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return "", 0, err
	}
	return host, uint16(portInt), nil
}
