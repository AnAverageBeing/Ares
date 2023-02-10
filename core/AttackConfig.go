package core

import (
	"net"
	"time"
)

type AttackConfig struct {
	Host         string
	Version      int
	ProxyManager *ProxyManager
	PerDelay     int
	Delay        time.Duration
	Loops        int
}

const DefaultPort = "25565"

func NewConfig(srvAddr string, version int, manager *ProxyManager, perDelay int, delay time.Duration, loops int) (cfg *AttackConfig) {
	addr, port, _ := net.SplitHostPort(srvAddr)
	if port == "" {
		port = DefaultPort
	}

	cfg = &AttackConfig{
		Host:         net.JoinHostPort(addr, port),
		Version:      version,
		ProxyManager: manager,
		PerDelay:     perDelay,
		Delay:        delay,
		Loops:        loops,
	}

	return
}
