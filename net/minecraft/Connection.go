package minecraft

import (
	pk "Ares/net/minecraft/packet"
	"Ares/net/proxy"
	"crypto/cipher"
	"io"
	"net"
)

type Connection struct {
	Socket net.Conn
	Reader io.Reader
	Writer io.Writer

	threshold int
}

func DialMc(host string, proxy *proxy.Proxy) (connection *Connection, err error) {
	conn, err := proxy.Dial(host)
	connection = WrapConn(*conn)
	return
}

func WrapConn(conn net.Conn) *Connection {
	return &Connection{
		Socket:    conn,
		Reader:    conn,
		Writer:    conn,
		threshold: -1,
	}
}

func (c *Connection) Close() error { return c.Socket.Close() }

func (c *Connection) ReadPacket(p *pk.Packet) error {
	return p.UnPack(c.Reader, c.threshold)
}

func (c *Connection) WritePacket(p pk.Packet) error {
	return p.Pack(c.Writer, c.threshold)
}

func (c *Connection) SetCipher(ecoStream, decoStream cipher.Stream) {
	c.Reader = cipher.StreamReader{ // Set receiver for AES
		S: decoStream,
		R: c.Socket,
	}
	c.Writer = cipher.StreamWriter{
		S: ecoStream,
		W: c.Socket,
	}
}

func (c *Connection) SetThreshold(t int) {
	c.threshold = t
}
