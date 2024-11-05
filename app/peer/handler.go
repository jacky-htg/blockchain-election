package peer

import (
	"bufio"
	"fmt"
	"myapp/app/blockchain"
	"net"
	"time"

	"github.com/bytedance/sonic"
)

type MessageType string

const (
	BlockchainUpdate  MessageType = "blockchain_update"
	RequestBlockchain MessageType = "request_blockchain"
	NewPeerJoined     MessageType = "new_peer_joined"
	ShutdownPeer      MessageType = "shutdown_peer"
)

type Message struct {
	Type MessageType `json:"type"`
	Data []byte      `json:"data"`
}

// Handler untuk setiap koneksi peer
func (p2p *P2PNetwork) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Baca request dari client
	data, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}

	var incomingMessage Message
	if err := sonic.Unmarshal(data, &incomingMessage); err != nil {
		fmt.Println("Error decoding blockchain:", err)
		return
	}

	switch incomingMessage.Type {
	case BlockchainUpdate:
		// Verifikasi dan update blockchain
		var blockchainData *blockchain.Blockchain
		if err := sonic.Unmarshal(incomingMessage.Data, &blockchainData); err != nil {
			fmt.Println("Error decoding blockchain:", err)
			return
		}
		if p2p.VerifyAndUpdateBlockchain(blockchainData) {
			fmt.Println("Updated local blockchain with incoming blockchain.")
		} else {
			fmt.Println("Received invalid or shorter blockchain.")
		}

	case RequestBlockchain:
		blockchainData, _ := sonic.Marshal(p2p.Blockchain)
		message := Message{Type: BlockchainUpdate, Data: blockchainData}
		response, _ := sonic.Marshal(message)
		conn.Write(response)

	case NewPeerJoined:
		p2p.AddPeer(string(incomingMessage.Data))

	case ShutdownPeer:
		p2p.RemovePeer(string(incomingMessage.Data))
		fmt.Println(p2p.peers)

	default:
		fmt.Println("Received unknown message type:", incomingMessage.Type)
	}
}

// Verifikasi dan update blockchain jika lebih panjang
func (p2p *P2PNetwork) VerifyAndUpdateBlockchain(incoming *blockchain.Blockchain) bool {
	if !incoming.IsValid() {
		return false
	}

	if len(incoming.Blocks) > len(p2p.Blockchain.Blocks) {
		p2p.Blockchain = incoming
		for _, b := range p2p.Blockchain.Blocks {
			fmt.Println(b.Data)
		}
		return true
	}

	return false
}

// Fungsi untuk menangani suara dari voter dan menambahkan blok baru.
func (p2p *P2PNetwork) HandleVote(voterID, candidateID string) {
	for _, block := range p2p.Blockchain.Blocks {
		if block.Data.VoterID == voterID {
			fmt.Println("Voter already voted:", voterID)
			return
		}
	}

	voteData := blockchain.VoteData{
		VoterID:     voterID,
		CandidateID: candidateID,
		Timestamp:   time.Now().Unix(),
	}

	newBlock := blockchain.Block{
		Index:     len(p2p.Blockchain.Blocks),
		Timestamp: time.Now().Unix(),
		Data:      voteData,
		PrevHash:  []byte{},
	}

	if len(p2p.Blockchain.Blocks) > 0 {
		newBlock.PrevHash = p2p.Blockchain.Blocks[len(p2p.Blockchain.Blocks)-1].Hash
	}

	pow := blockchain.NewProofOfWork(&newBlock)
	nonce, hash := pow.Run()
	newBlock.Hash = hash
	newBlock.Nonce = nonce

	if pow.Validate() {
		if p2p.Blockchain.AddBlock(newBlock) {
			err := p2p.Blockchain.Election.Vote(voterID, candidateID) // Pastikan Election ada di blockchain
			if err != nil {
				fmt.Println("Error while voting:", err)
			} else {
				p2p.BroadcastBlockchain()
				fmt.Println("Vote successful for voter:", voterID)
			}
		}
	} else {
		fmt.Println("Failed to validate block")
	}
}

func (p2p *P2PNetwork) RequestBlockchainFromPeers() {
	for _, peer := range p2p.peers {
		conn, err := net.Dial("tcp", peer.Address)
		if err != nil {
			fmt.Println("Error connecting to peer:", err)
			continue
		}
		defer conn.Close()

		// Kirim permintaan untuk mendapatkan blockchain
		message := Message{Type: RequestBlockchain}
		data, _ := sonic.Marshal(message)
		writer := bufio.NewWriter(conn)
		data = append(data, '\n')

		_, err = writer.WriteString(string(data))
		if err != nil {
			fmt.Println("gagal mengirim request:", err)
			continue
		}
		writer.Flush()

		// Terima blockchain dari peer
		reader := bufio.NewReader(conn)
		peerData, _ := reader.ReadBytes('\n')
		var receivedMessage Message
		if err := sonic.Unmarshal(peerData, &receivedMessage); err != nil {
			fmt.Println("Error decoding blockchain from peer:", err)
			continue
		}

		// Pastikan type data adalah BlockchainUpdate sebelum disinkronkan
		if receivedMessage.Type == BlockchainUpdate {
			var peerBlockchain *blockchain.Blockchain
			if err := sonic.Unmarshal(receivedMessage.Data, &peerBlockchain); err != nil {
				fmt.Println("Error decoding blockchain blocks:", err)
				continue
			}
			p2p.Blockchain.SyncWithPeer(peerBlockchain.Blocks, peerBlockchain.Election)
			fmt.Println("Synchronized blockchain with peer:", peer.Address)
			break // Stop setelah sinkronisasi dengan satu peer
		}
	}
}
