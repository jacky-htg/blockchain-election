package peer

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/bytedance/sonic"
)

type BaseRequest struct {
	Type string `json:"type"`
}

func (p2p *P2PNetwork) ListenForBlocks(address string) {
	ln, err := net.Listen("tcp", address)
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

func (p2p *P2PNetwork) GetPeersFromBootstrap() ([]Peer, error) {
	conn, err := net.DialTimeout("tcp", p2p.bootstrap, 30*time.Second)
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke bootstrap server: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	req := BaseRequest{Type: "GET_PEERS"}
	data, err := sonic.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat request JSON: %v", err)
	}
	writer := bufio.NewWriter(conn)
	data = append(data, '\n')

	_, err = writer.WriteString(string(data))
	if err != nil {
		return nil, fmt.Errorf("gagal mengirim request: %v", err)
	}
	writer.Flush()

	// Membaca response dari server
	response, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("gagal membaca response: %v", err)
	}

	var peers []Peer

	err = sonic.Unmarshal([]byte(strings.TrimSpace(response)), &peers)
	if err != nil {
		return nil, fmt.Errorf("gagal melakukan unmarshal response: %v", err)
	}

	return peers, nil
}
