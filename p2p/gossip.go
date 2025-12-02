package p2p

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"github.com/OhMyDitzzy/vulcan/core"
	"github.com/OhMyDitzzy/vulcan/types"
	"github.com/OhMyDitzzy/vulcan/txpool"
)

type Node struct {
	port       int
	peers      []*Peer
	blockchain *core.Blockchain
	mempool    *txpool.Mempool
	listener   net.Listener
	mu         sync.RWMutex
	running    bool
}

func NewNode(port int, bc *core.Blockchain, mp *txpool.Mempool, bootstrapPeers []string) *Node {
	node := &Node{
		port:       port,
		blockchain: bc,
		mempool:    mp,
		peers:      make([]*Peer, 0),
	}
	
	// Connect to bootstrap peers
	for _, addr := range bootstrapPeers {
		peer := NewPeer(addr)
		if err := peer.Connect(); err != nil {
			log.Printf("Failed to connect to peer %s: %v", addr, err)
		} else {
			node.peers = append(node.peers, peer)
		}
	}
	
	return node
}

func (n *Node) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", n.port))
	if err != nil {
		return err
	}
	
	n.listener = listener
	n.running = true
	
	go n.acceptConnections()
	return nil
}

func (n *Node) Stop() {
	n.running = false
	if n.listener != nil {
		n.listener.Close()
	}
	
	for _, peer := range n.peers {
		peer.Close()
	}
}

func (n *Node) acceptConnections() {
	for n.running {
		conn, err := n.listener.Accept()
		if err != nil {
			if n.running {
				log.Printf("Accept error: %v", err)
			}
			continue
		}
		
		go n.handleConnection(conn)
	}
}

func (n *Node) handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	
	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}
		
		n.handleMessage(&msg)
	}
}

func (n *Node) handleMessage(msg *Message) {
	switch msg.Type {
	case "new_transaction":
		var tx types.Transaction
		if err := json.Unmarshal(msg.Data, &tx); err == nil {
			n.mempool.AddTransaction(&tx)
			n.BroadcastTransaction(&tx)
		}
	case "new_block":
		var block core.Block
		if err := json.Unmarshal(msg.Data, &block); err == nil {
			n.blockchain.AddBlock(&block)
			n.BroadcastBlock(&block)
		}
	}
}

func (n *Node) BroadcastTransaction(tx *types.Transaction) {
	data, _ := json.Marshal(tx)
	msg := &Message{Type: "new_transaction", Data: data}
	
	for _, peer := range n.peers {
		peer.SendMessage(msg)
	}
}

func (n *Node) BroadcastBlock(block *core.Block) {
	data, _ := json.Marshal(block)
	msg := &Message{Type: "new_block", Data: data}
	
	for _, peer := range n.peers {
		peer.SendMessage(msg)
	}
}

func (n *Node) GetPeers() []string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	
	peers := make([]string, len(n.peers))
	for i, p := range n.peers {
		peers[i] = p.Address
	}
	return peers
}

func (n *Node) AddPeer(address string) error {
	peer := NewPeer(address)
	if err := peer.Connect(); err != nil {
		return err
	}
	
	n.mu.Lock()
	n.peers = append(n.peers, peer)
	n.mu.Unlock()
	
	return nil
}