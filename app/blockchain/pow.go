package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/bytedance/sonic"
)

const targetBits = 24

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) prepareData(nonce int) ([]byte, error) {
	dataBytes, err := sonic.Marshal(pow.block.Data)
	if err != nil {
		return nil, err // Kembalikan error jika terjadi kesalahan
	}
	data := bytes.Join(
		[][]byte{
			[]byte(fmt.Sprintf("%d", pow.block.Index)),
			[]byte(fmt.Sprintf("%d", pow.block.Timestamp)),
			dataBytes,
			pow.block.PrevHash,
			[]byte(fmt.Sprintf("%d", nonce)),
		},
		[]byte{},
	)
	return data, nil
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for nonce < (1<<63 - 1) {
		data, err := pow.prepareData(nonce)
		if err != nil {
			fmt.Println("Error preparing data:", err)
			return -1, nil
		}
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data, err := pow.prepareData(pow.block.Nonce)
	if err != nil {
		fmt.Println("Error preparing data:", err)
		return false
	}
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.target) == -1
}
