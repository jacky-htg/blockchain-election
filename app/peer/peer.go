package peer

import (
	"bufio"
	"crypto/ecdsa"
	"fmt"
	"myapp/app/blockchain"
	"net"

	"github.com/bytedance/sonic"
)

type RegisterRequest struct {
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
}

// Peer struct
type Peer struct {
	Address   string `json:"address"`
	PublicKey string `json:"publicKey"`
}

type P2PNetwork struct {
	peers      []Peer
	Blockchain *blockchain.Blockchain
	bootstrap  string // Alamat bootstrap server
	localAddr  string // Alamat lokal peer
	privateKey *ecdsa.PrivateKey
	publicKey  []byte
}

func NewP2PNetwork(bootstrap, localAddr string, privateKey *ecdsa.PrivateKey, publicKey []byte) *P2PNetwork {
	return &P2PNetwork{
		peers:      make([]Peer, 0),
		Blockchain: &blockchain.Blockchain{},
		bootstrap:  bootstrap,
		localAddr:  localAddr,
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (p2p *P2PNetwork) AddPeer(newPeer Peer) {
	for _, peer := range p2p.peers {
		if peer.Address == newPeer.Address {
			return // Peer sudah ada
		}
	}
	p2p.peers = append(p2p.peers, newPeer)
	fmt.Println(p2p.peers)
}

func (p2p *P2PNetwork) RemovePeer(address string) {
	for i, peer := range p2p.peers {
		if peer.Address == address {
			p2p.peers = append(p2p.peers[:i], p2p.peers[i+1:]...)
			return
		}
	}
	fmt.Println(p2p.peers)
}

func (p2p *P2PNetwork) RegisterToBootstrap() error {
	conn, err := net.Dial("tcp", p2p.bootstrap)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke bootstrap server: %v", err)
	}
	defer conn.Close()

	var peer Peer = Peer{
		Address:   p2p.localAddr,
		PublicKey: string(p2p.publicKey),
	}

	payload, err := sonic.Marshal(peer)
	if err != nil {
		return fmt.Errorf("gagal melakukan marshal data: %v", err)
	}

	request := RegisterRequest{
		Type:    "REGISTER",
		Payload: payload,
	}
	data, err := sonic.Marshal(request)
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

	var peer Peer = Peer{
		Address: p2p.localAddr,
	}

	payload, _ := sonic.Marshal(peer)

	request := RegisterRequest{
		Type:    "REMOVE",
		Payload: payload,
	}
	data, _ := sonic.Marshal(request)
	conn.Write(append(data, '\n'))
}
