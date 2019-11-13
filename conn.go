package gstats

import (
	"net"
)

type conn struct {
	net.Conn
	g *GStats
}

func (c *conn) Read(b []byte) (int, error) {
	n, err := c.Conn.Read(b)
	c.g.notifyConnRead(c.Conn.RemoteAddr(), n)
	return n, err
}

func (c *conn) Write(b []byte) (int, error) {
	n, err := c.Conn.Write(b)
	c.g.notifyConnWrite(c.Conn.RemoteAddr(), n)
	return n, err
}

func (c *conn) Close() error {
	c.g.notifyConnClose(c.Conn.RemoteAddr())
	return c.Conn.Close()
}
