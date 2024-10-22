package main

import (
	"fmt"
	"time"

	"myapp/app/blockchain"
)

func main() {
	// Membuat blockchain baru
	bc := blockchain.Blockchain{}

	// Menambahkan genesis block
	bc.AddBlock("Genesis Block")

	// Menambahkan beberapa transaksi
	bc.AddBlock("Transaction 1")
	bc.AddBlock("Transaction 2")

	// Menampilkan semua blok dalam blockchain
	fmt.Println("Blockchain:")
	for _, block := range bc.Blocks {
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Timestamp: %s\n", time.Unix(block.Timestamp, 0))
		fmt.Printf("Data: %s\n", block.Data.Data)
		fmt.Printf("PrevHash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println("------------------------------")
	}

	// Validasi blockchain
	fmt.Println("Is blockchain valid?", bc.IsValid())
}
