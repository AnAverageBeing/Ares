package utils

import (
	"Ares/core"
	"Ares/net/proxy"
	"bufio"
	"os"
	"strings"
	"time"
)

func LoadFromFile(protocol proxy.Protocol, timeout time.Duration, path string, manager *core.ProxyManager) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		proxy := &proxy.Proxy{
			Host:     line,
			Protocol: protocol,
			Timeout:  timeout,
		}

		manager.Add(proxy)
	}

	return scanner.Err()
}
