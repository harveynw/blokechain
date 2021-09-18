package miner

import (
	"fmt"
	"time"
	"github.com/harveynw/blokechain/internal/chain"
	"github.com/harveynw/blokechain/internal/data"
	"github.com/harveynw/blokechain/internal/params"
)

// Mine performs a fixed number of iterations attempting to mine a block header
func Mine(bh chain.BlockHeader) (success bool, bhSolved chain.BlockHeader) {
	target := bh.DifficultyTarget

	// // For comparison
	// lowest := make([]byte, 32)
	// for i := range lowest {
	// 	lowest[i] = 0xFF
	// }


	for i := 0; i < params.MiningIterationsPerCall; i++ {
		work := data.DoubleHash(bh.Encode(), false)

		// if chain.Compare(work, lowest) == -1 {
		// 	lowest = work
		// }

		if target.IsSolution(work) {
			return true, bh
		}

		bh = bh.IncrementNonce()
	}

	// lowestFound := new(big.Int)
	// lowestFound.SetBytes(lowest)
	// fmt.Println("Lowest", lowestFound)

	return false, bh
}

// MinerTest tests
func MinerTest() {
	for i := 1; i <= 256; i+=2 {
		genBlockHeader := chain.Genesis(int64(i)).Header
		start := time.Now()
		success := false
		for !success {
			success, *genBlockHeader = Mine(*genBlockHeader)
			//fmt.Println(params.MiningIterationsPerCall, "hashing iterations completed")
		}
		elapsed := time.Since(start)
		fmt.Printf("%v \n", elapsed.Milliseconds())
	}
	// genBlockHeader := chain.Genesis().Header

	// success := false
	// //var solvedBlockHeader chain.BlockHeader

	// start := time.Now()
	// for !success {
	// 	success, _ = Mine(*genBlockHeader)
	// 	fmt.Println(params.MiningIterationsPerCall, "hashing iterations completed")
	// }
	// elapsed := time.Since(start)
    // fmt.Printf("Mine took %s \n", elapsed)

	// //fmt.Println(success)
	// //fmt.Println(solvedBlockHeader)
	// //fmt.Println("Nonce is", solvedBlockHeader.Nonce)
}