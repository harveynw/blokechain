package peer

import (
	"fmt"
)

type PeerList struct {
	peers []Peer
}

func NewPeerList() *PeerList {
	peers := make([]Peer, 0)
	return &PeerList{peers}
}

func (p *PeerList) AddPeer(peer *Peer) {
	fmt.Println("Adding Peer:", peer)
	p.peers = append(p.peers, *peer)
}
