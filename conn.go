package gstats

import (
	"net"
	"sync"
)

type conn struct {
	net.Conn
	g *GStats
	o *sync.Once
}

func (c *conn) Close() error {
	c.o.Do(func() {
		c.g.notifyConnClose(c.Conn.RemoteAddr())
	})
	return c.Conn.Close()
}
