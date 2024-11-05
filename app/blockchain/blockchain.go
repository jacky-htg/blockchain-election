package blockchain

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

type Blockchain struct {
	Blocks   []Block
	Election *Election
	mu       sync.Mutex
}

func (bc *Blockchain) SetGenesisBlock() bool {
	if len(bc.Blocks) == 0 {
		// Membuat dan membroadcast blok genesis.
		genesisBlock := Block{
			Index:     0,
			Timestamp: time.Now().Unix(),
			Data:      VoteData{VoterID: "system", CandidateID: "Genesis Block"},
		}
		pow := NewProofOfWork(&genesisBlock)
		nonce, hash := pow.Run()
		genesisBlock.Hash = hash
		genesisBlock.Nonce = nonce

		if pow.Validate() {
			if bc.AddBlock(genesisBlock) {
				fmt.Println("Added genesis block:", genesisBlock.Data.VoterID)
				return true
			}
		} else {
			fmt.Println("Failed to validate genesis block")
		}
	}

	return false
}

func (bc *Blockchain) AddBlock(newBlock Block) bool {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	for _, block := range bc.Blocks {
		if bytes.Equal(block.Hash, newBlock.Hash) {
			return false
		}
	}
	bc.Blocks = append(bc.Blocks, newBlock)
	return true
}

func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		if !bytes.Equal(currentBlock.PrevHash, prevBlock.Hash) {
			return false
		}

		pow := NewProofOfWork(&currentBlock)
		if !pow.Validate() {
			return false
		}
	}
	return true
}

func (bc *Blockchain) SyncWithPeer(peerBlocks []Block, election *Election) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if len(peerBlocks) > len(bc.Blocks) {
		bc.Blocks = peerBlocks
		bc.Election = election
	}
}
