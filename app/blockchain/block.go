package blockchain

type Block struct {
	Index     int
	Timestamp int64
	Data      VoteData
	PrevHash  []byte
	Hash      []byte
	Nonce     int
}
