package wallet

import (
	"fmt"
)

func Test() {
	wallet := Load()
	wallet.GenerateNew()
	wallet.GenerateNew()

	fmt.Println(wallet.ListAddresses())
}