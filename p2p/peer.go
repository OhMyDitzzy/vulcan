package p2p

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type Peer struct {
	Address string
	conn    net.Conn
	mu      sync.Mutex
}

func NewPeer(address string) *Peer {
	return &Peer{Address: address}
}

func (p *Peer) Connect() error {
	conn, err := net.Dial("tcp", p.Address)
	if err != nil {
		return err
	}
	p.conn = conn
	return nil
}

func (p *Peer) SendMessage(msg *Message) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.conn == nil {
		return fmt.Errorf("not connected")
	}
	
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	
	_, err = p.conn.Write(append(data, '\n'))
	return err
}

func (p *Peer) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}