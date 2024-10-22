package peer

import (
	"encoding/json"
	"fmt"
	"myapp/app/blockchain"
	"net"
)

type Peer struct {
	address string
}

type P2PNetwork struct {
	peers      []Peer
	Blockchain *blockchain.Blockchain
}

func NewP2PNetwork() *P2PNetwork {
	return &P2PNetwork{
		peers:      make([]Peer, 0),
		Blockchain: &blockchain.Blockchain{},
	}
}

func (p2p *P2PNetwork) AddPeer(address string) {
	peer := Peer{address: address}
	p2p.peers = append(p2p.peers, peer)
}

func (p2p *P2PNetwork) BroadcastBlock(block blockchain.Block) {
	for _, peer := range p2p.peers {
		conn, err := net.Dial("tcp", peer.address)
		if err != nil {
			fmt.Println("Error connecting to peer:", err)
			continue
		}
		defer conn.Close()

		encoder := json.NewEncoder(conn)
		if err := encoder.Encode(block); err != nil {
			fmt.Println("Error encoding block:", err)
		}
	}
}

func (p2p *P2PNetwork) ListenForBlocks(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error setting up listener:", err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go p2p.handleConnection(conn)
	}
}

func (p2p *P2PNetwork) handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	var block blockchain.Block
	if err := decoder.Decode(&block); err != nil {
		fmt.Println("Error decoding block:", err)
		return
	}

	p2p.Blockchain.AddBlock(block.Data.Data)
	fmt.Println("Received new block:", block)
}
