package peer

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Peer struct {
	conn *net.TCPConn
	w *bufio.Writer
}

func newPeer(conn *net.TCPConn) *Peer {
	p := Peer{conn, bufio.NewWriter(conn)}

	fmt.Println("New Peer established", conn.LocalAddr(), conn.RemoteAddr())

	go p.loopCheckAlive()

	return &p
}

func (p *Peer) loopCheckAlive() {
	for {
		fmt.Fprintf(p.conn, "PING")
		time.Sleep(1000)
	}
}

