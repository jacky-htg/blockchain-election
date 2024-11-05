package peer

import (
	"bufio"
	"fmt"
	"myapp/app/blockchain"
	"net"

	"github.com/bytedance/sonic"
)

type RegisterRequest struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

// Peer struct
type Peer struct {
	Address string `json:"address"`
}

type P2PNetwork struct {
	peers      []Peer
	Blockchain *blockchain.Blockchain
	bootstrap  string // Alamat bootstrap server
	localAddr  string // Alamat lokal peer
}

func NewP2PNetwork(bootstrap, localAddr string) *P2PNetwork {
	return &P2PNetwork{
		peers:      make([]Peer, 0),
		Blockchain: &blockchain.Blockchain{},
		bootstrap:  bootstrap,
		localAddr:  localAddr,
	}
}

func (p2p *P2PNetwork) AddPeer(address string) {
	for _, peer := range p2p.peers {
		if peer.Address == address {
			return // Peer sudah ada
		}
	}
	peer := Peer{Address: address}
	p2p.peers = append(p2p.peers, peer)
}

func (p2p *P2PNetwork) RemovePeer(address string) {
	for i, peer := range p2p.peers {
		if peer.Address == address {
			p2p.peers = append(p2p.peers[:i], p2p.peers[i+1:]...)
			return
		}
	}
}

func (p2p *P2PNetwork) RegisterToBootstrap() error {
	conn, err := net.Dial("tcp", p2p.bootstrap)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke bootstrap server: %v", err)
	}
	defer conn.Close()

	// Membuat payload untuk didaftarkan ke bootstrap
	payload := RegisterRequest{
		Type:    "REGISTER",
		Payload: p2p.localAddr,
	}
	data, err := sonic.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gagal melakukan marshal data: %v", err)
	}

	writer := bufio.NewWriter(conn)
	data = append(data, '\n')

	_, err = writer.WriteString(string(data))
	if err != nil {
		return fmt.Errorf("gagal mengirim request: %v", err)
	}
	writer.Flush()

	fmt.Println("Berhasil mendaftar ke bootstrap server")
	return nil
}

func (p2p *P2PNetwork) NotifyBootstrapOnShutdown() {
	conn, err := net.Dial("tcp", p2p.bootstrap)
	if err != nil {
		fmt.Println("Error connecting to bootstrap server:", err)
		return
	}
	defer conn.Close()

	payload := RegisterRequest{
		Type:    "REMOVE",
		Payload: p2p.localAddr,
	}
	data, _ := sonic.Marshal(payload)
	conn.Write(append(data, '\n'))
}
