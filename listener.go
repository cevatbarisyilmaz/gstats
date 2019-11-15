package gstats

import (
	"net"
	"sync"
)

type listener struct {
	net.Listener
	g *GStats
}

func (l *listener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	l.g.notifyNewConn(c.RemoteAddr())
	return &conn{
		Conn: c,
		g:    l.g,
		o:    &sync.Once{},
	}, nil
}
