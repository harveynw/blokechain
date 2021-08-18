package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"github.com/harveynw/blokechain/internal/wallet"
)

// ErrTooManyArguments when overspecified
var ErrTooManyArguments = errors.New("Too many arguments")

// Option defines a wallet CLI command
type Option struct {
	handler func()
	triggered *bool
	description string
}

func newOption(command string, handler func(), description string) Option {
	f := flag.Bool(command, false, description)
	return Option{handler: handler, triggered: f, description: description}
}

var options = []Option {
	newOption("g", newAddress, "Generate a new address"),
	newOption("ls", listAddresses, "List addresses"),
	newOption("b", listBalances, "List balances of addresses"),
}

func main() {
	flag.Parse()
	assertSingleArgument()

	for _, option := range options {
		if *option.triggered {
			option.handler()
			break
		}
	}
}

func assertSingleArgument() {
	if flag.NArg() + flag.NFlag() != 1 {
		fmt.Println(ErrTooManyArguments)
		os.Exit(0)
	}
}

func newAddress() {
	assertSingleArgument()

	addr := wallet.Load().GenerateNew()
	fmt.Println(addr)
}

func listAddresses() {
	assertSingleArgument()

	addresses := wallet.Load().Addresses
	for _, addr := range addresses {
		fmt.Println(addr.PublicKey.ToAddress())
	}
}

func listBalances() {
	assertSingleArgument()

	fmt.Printf("Gathering UXTOs from nodes...\n\n")
	fmt.Println("Total Balance : 0 BTC")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "\t\033[1mAddress\033[0m\t\033[1mBalance\033[0m\t")

	addresses := wallet.Load().Addresses
	for _, addr := range addresses {
		fmt.Fprintln(w, "\t\033[0m" + addr.PublicKey.ToAddress() + "\033[0m\t\033[0m0\033[0m\t")
	}
	w.Flush()
}