package blockchain

import (
	"bytes"
	"time"
)

type Blockchain struct {
	Blocks []Block
}

func (bc *Blockchain) AddBlock(data string) Block {
	var newBlock Block
	if len(bc.Blocks) == 0 {
		newBlock = createBlock(0, "Genesis Block", []byte{})
	} else {
		prevBlock := bc.Blocks[len(bc.Blocks)-1]
		newBlock = createBlock(prevBlock.Index+1, data, prevBlock.Hash)
	}
	bc.Blocks = append(bc.Blocks, newBlock)
	return newBlock
}

func createBlock(index int, data string, prevHash []byte) Block {
	newData := Data{Data: data}
	block := Block{
		Index:     index,
		Timestamp: time.Now().Unix(),
		Data:      newData,
		PrevHash:  prevHash,
		Hash:      []byte{},
	}
	block.Hash = block.calculateHash()
	return block
}

func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		if !bytes.Equal(currentBlock.Hash, currentBlock.calculateHash()) {
			return false
		}
		if !bytes.Equal(currentBlock.PrevHash, prevBlock.Hash) {
			return false
		}
	}
	return true
}
