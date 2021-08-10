// package main

// import (
// 	"bufio"
// 	"crypto/sha512"
// 	"crypto/x509"
// 	"flag"
// 	"fmt"
// 	"net"
// 	"github.com/harveynw/blokechain/internal/data"
// 	"github.com/harveynw/blokechain/internal/peer"
// )

// func main() {
// 	privateKey := data.GenerateValidPrivateKey()
// 	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println("Private Key:", string(privateKeyBytes))


// 	fmt.Println("Starting blokechain node...")

// 	startGenesisBlock := flag.Bool("genesis", false, "Flag that triggers starting a new blockchain")

// 	startServer := flag.Bool("server", false, "Start test server")
// 	address := flag.String("address", ":8080", "Address to connect to/start server on")

// 	flag.Parse()

// 	if *startGenesisBlock {
// 		fmt.Println("Beginning Genesis Block!")
// 	} else {
// 		fmt.Println("Attempting to connect to peers")
// 	}

// 	data := []byte("hllo")
// 	hash := sha512.Sum512(data)
// 	fmt.Printf("%x\n", hash)

// 	peers := peer.NewPeerList()

// 	if *startServer {
// 		startServerListen(*address, peers)
// 	} else {
// 		startClientConnect(*address, peers)
// 	}
// }

// func startServerListen(address string) {
// 	ln, err := net.Listen("tcp", address)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println("Listening")
// 	for {
// 		conn, err := ln.Accept()
// 		if err != nil {
// 			fmt.Println(err)
// 		} else {
// 			fmt.Println("Connection!", conn)
// 		}
// 	}

// }

// func startClientConnect(address string) {
// 	conn, err := net.Dial("tcp", address)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")

// 	status, _ := bufio.NewReader(conn).ReadString('\n')

// 	fmt.Println("Status", status)

// }