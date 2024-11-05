package peer

import (
	"bufio"
	"fmt"
	"net"

	"github.com/bytedance/sonic"
)

// BroadcastBlockchain mengirimkan seluruh blockchain ke semua peers yang terhubung.
func (p2p *P2PNetwork) BroadcastBlockchain() {
	for _, peer := range p2p.peers {
		conn, err := net.Dial("tcp", peer.Address)
		if err != nil {
			fmt.Println("Error connecting to peer:", err)
			continue
		}
		defer conn.Close()

		blockchainData, err := sonic.Marshal(p2p.Blockchain)
		if err != nil {
			fmt.Println("Gagal membuat JSON untuk Blockchain:", err)
			continue
		}
		message := Message{
			Type: BlockchainUpdate,
			Data: blockchainData,
		}
		data, err := sonic.Marshal(message)
		if err != nil {
			fmt.Println("gagal membuat request JSON:", err)
			continue
		}
		writer := bufio.NewWriter(conn)
		data = append(data, '\n')

		_, err = writer.WriteString(string(data))
		if err != nil {
			fmt.Println("gagal mengirim request:", err)
			continue
		}
		writer.Flush()
	}
}
