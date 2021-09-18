package main

import (
	"github.com/harveynw/blokechain/internal/miner"
	// "github.com/harveynw/blokechain/internal/script"
	// "github.com/harveynw/blokechain/internal/wallet"
)

func main() {
	//data.TestEncoding()
	// pkEncode := wallet.Load().Addresses[0].PublicKey.HashEncode()

	// s := script.P2PKH(pkEncode).Encode()

	// for _, v := range s {
	// 	fmt.Printf("%v,", v)
	// }
	// fmt.Println("")

	miner.MinerTest()
}

