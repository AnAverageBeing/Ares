package proxy

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

func (p Proxy) dialHTTP(target string) (*net.Conn, error) {
	dialer := &net.Dialer{Timeout: p.Timeout}
	proxyURL := &url.URL{
		Scheme: "http",
		Host:   p.Host,
	}
	httpTransport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		},
	}
	client := &http.Client{Transport: httpTransport}

	conn, err := dialer.Dial("tcp", target)
	if err != nil {
		return nil, err
	}

	if p.Auth != nil {
		auth := fmt.Sprintf("%s:%s", p.Auth.Username, p.Auth.Password)
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
		req, err := http.NewRequest("CONNECT", target, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Proxy-Authorization", fmt.Sprintf("Basic %s", encodedAuth))
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("proxy authentication failed: %s", resp.Status)
		}
	}

	return &conn, nil
}
