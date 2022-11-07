package main

import (
	// "github.com/harveynw/blokechain/internal/miner"
	"github.com/harveynw/blokechain/internal/script"
	// "github.com/harveynw/blokechain/internal/wallet"
	// "github.com/harveynw/blokechain/internal/cryptography"
)

func main() {
	//data.TestEncoding()
	// pkEncode := wallet.Load().Addresses[0].PublicKey.HashEncode()

	// s := script.P2PKH(pkEncode).Encode()

	// for _, v := range s {
	// 	fmt.Printf("%v,", v)
	// }
	// fmt.Println("")

	//miner.MinerTest()

	// cryptography.RIPEMD160([]byte{0x00, 0x05, 0xac})

	// Testing slice pointers
	// stack := []byte{0x01, 0x02, 0x03, 0x04}
	// stackref := &stack
	// *stackref = (*stackref)[:len(*stackref)-1]
	// fmt.Println(stack)
	script.TestScript()
}

