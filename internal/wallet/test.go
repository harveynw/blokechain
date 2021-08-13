package wallet

import (
	"fmt"
)

func Test() {
	wallet := Load()
	for i := 0; i < 5; i++ {
		wallet.GenerateNew()
	}
	fmt.Println(wallet.ListAddresses())
}