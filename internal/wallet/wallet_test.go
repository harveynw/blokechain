package wallet

import (
	"fmt"
	"os"
	"testing"
)

func TestLoading(t *testing.T) {
	os.Chdir("../..")

	wallet := Load()
	for i := 0; i < 5; i++ {
		wallet.GenerateNew()
	}
	fmt.Println(wallet.ListAddresses())
}