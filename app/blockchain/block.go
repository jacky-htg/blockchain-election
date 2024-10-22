package blockchain

import (
	"crypto/sha256"
	"fmt"
)

type Block struct {
	Index     int
	Timestamp int64
	Data      Data
	PrevHash  []byte
	Hash      []byte
}

func (b *Block) calculateHash() []byte {
	record := fmt.Sprintf("%d%d%s%s", b.Index, b.Timestamp, b.Data.Data, b.PrevHash)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hashed
}
