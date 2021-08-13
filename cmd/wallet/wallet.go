package main

import (
	"fmt"
	"flag"
	"errors"
	"github.com/harveynw/blokechain/internal/wallet"
)

// ErrTooManyArguments when overspecified
var ErrTooManyArguments = errors.New("Too many arguments")

func main() {
	generateNew := flag.Bool("g", false, "Generate a new address")
	listAll := flag.Bool("ls", false, "List addresses")

	flag.Parse()

	if *generateNew {
		newAddress()
	} else if *listAll {
		listAddresses()
	}
}

func newAddress() {
	if flag.NArg() + flag.NFlag() > 1 {
		fmt.Println(ErrTooManyArguments)
		return
	}

	addr := wallet.Load().GenerateNew()
	fmt.Println(addr)
}

func listAddresses() {
	if flag.NArg() + flag.NFlag() > 1 {
		fmt.Println(ErrTooManyArguments)
		return
	}

	addresses := wallet.Load().Addresses
	for _, addr := range addresses {
		fmt.Println(addr.PublicKey.ToAddress())
	}
}