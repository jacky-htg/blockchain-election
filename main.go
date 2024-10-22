package main

import (
	"bufio"
	"flag"
	"fmt"
	"myapp/app/blockchain"
	"myapp/app/peer"
	"os"
)

func readInput(p2p *peer.P2PNetwork) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data := scanner.Text()
		block := p2p.Blockchain.AddBlock(data)
		p2p.BroadcastBlock(block)
		fmt.Println("Added new block:", block)
	}
}

func main() {
	port := flag.String("port", "3000", "Port to listen on")
	flag.Parse()

	// Membuat jaringan P2P baru
	p2p := peer.NewP2PNetwork()

	// Menambahkan peer secara manual
	if *port == "3000" {
		p2p.AddPeer("localhost:3001")
		p2p.AddPeer("localhost:3002")
	} else if *port == "3001" {
		p2p.AddPeer("localhost:3000")
		p2p.AddPeer("localhost:3002")
	} else if *port == "3002" {
		p2p.AddPeer("localhost:3000")
		p2p.AddPeer("localhost:3001")
	}

	// Mendengarkan koneksi untuk menerima blok
	go p2p.ListenForBlocks(*port)

	if *port == "3000" {
		// Membroadcast blok genesis ke semua peer
		genesisBlock := blockchain.Block{Index: 0, Timestamp: 0, Data: blockchain.Data{Data: "Genesis Block"}}
		p2p.BroadcastBlock(genesisBlock)
	}

	// Membaca input dari pengguna untuk menambahkan blok baru dan membroadcast-nya
	go readInput(p2p)

	// Menunggu agar program tetap berjalan
	select {}
}
