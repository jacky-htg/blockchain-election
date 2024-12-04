package main

import (
	"bufio"
	"flag"
	"fmt"
	"myapp/app/blockchain"
	"myapp/app/peer"
	"myapp/app/pkg/signature"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	address := flag.String("address", "localhost:5000", "Address for node p2p network")
	init := flag.Bool("init", false, "init blockchain")
	flag.Parse()

	bootstrapAddress := "localhost:4000"
	// Membuat jaringan P2P dan kontrak voting.
	privateKey, err := signature.LoadOrCreateKeyPair("key.pem")
	if err != nil {
		fmt.Println("Error generating keys:", err)
		return
	}
	publicKey, err := signature.SerializePublicKey(&privateKey.PublicKey)
	if err != nil {
		fmt.Println("Error serializing public key:", err)
		return
	}
	println(string(publicKey))
	p2p := peer.NewP2PNetwork(bootstrapAddress, *address, privateKey, publicKey)
	p2p.RegisterToBootstrap()

	// Inisialisasi blockchain dengan instance Election
	p2p.Blockchain = &blockchain.Blockchain{
		Blocks:   []blockchain.Block{},
		Election: blockchain.NewElection([]string{}),
	}

	peers, err := p2p.GetPeersFromBootstrap()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, peer := range peers {
		fmt.Println("get peer from bootstrap", peer, peer.Address)
		if peer.Address != *address {
			p2p.AddPeer(peer)
			println("menambahkan: ", peer.Address)
		}
	}

	if *init {
		p2p.Blockchain.Election.AddCandidate("Alice")
		p2p.Blockchain.Election.AddCandidate("Bob")
		p2p.Blockchain.Election.AddCandidate("Charlie")
		println("prepare set genesis block")
		if !p2p.Blockchain.SetGenesisBlock() {
			p2p.BroadcastBlockchain()
		}
	} else {
		// Sinkronisasi blockchain untuk peer baru
		p2p.RequestBlockchainFromPeers()
	}

	// Mendengarkan koneksi untuk menerima blok.
	go p2p.ListenForBlocks(*address)

	go handleUserInput(p2p)

	// handling peer shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		p2p.NotifyBootstrapOnShutdown()
		os.Exit(0)
	}()

	// Menjaga agar program tetap berjalan.
	select {}
}

func handleUserInput(p2p *peer.P2PNetwork) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Ketik perintah. Contoh: vote voterID kandidatID atau showresult")
	for scanner.Scan() {
		input := scanner.Text()
		args := strings.Fields(input)

		if len(args) == 0 {
			fmt.Println("Masukkan perintah yang valid.")
			continue
		}

		switch args[0] {
		case "vote":
			if len(args) < 3 {
				fmt.Println("Perintah vote harus diikuti oleh voterID dan kandidatID.")
				continue
			}
			voterID := args[1]
			candidateID := args[2]
			p2p.HandleVote(voterID, candidateID)
			fmt.Printf("Vote dari %s untuk %s telah dicatat.\n", voterID, candidateID)

		case "showresult":
			fmt.Println("Hasil voting saat ini:")
			p2p.Blockchain.Election.DisplayResults()

		default:
			fmt.Println("Perintah tidak dikenal:", args[0])
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error membaca input:", err)
	}
}
