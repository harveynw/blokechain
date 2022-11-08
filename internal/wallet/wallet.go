package wallet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"github.com/harveynw/blokechain/internal/cryptography"
)

// Wallet for holding multiple key pairs
type Wallet struct {
	Addresses []Address
}

// Address for a holding a ECDSA public key, private key pair
type Address struct {
	PublicKey cryptography.PublicKey
	SecretKey *big.Int
}

// Save marshals wallet to file
func (wallet *Wallet) Save() {
	data, marshalErr := json.MarshalIndent(wallet, "", "   ")
	if marshalErr != nil {
		panic(marshalErr)
	}

	writeErr := ioutil.WriteFile(getWalletFile(), data, 0777)
	if writeErr != nil {
		panic(writeErr)
	}
}

// Load unmarshals the Wallet after retrieval from file
func Load() *Wallet {
	data, err := getWalletData()
	
	if err != nil {
		w := &Wallet{Addresses: make([]Address, 0)}
		w.Save()

		return w
	}

	w := &Wallet{}
	marshalErr := json.Unmarshal(data, w)

	if marshalErr != nil {
		panic(fmt.Errorf("%v", marshalErr))
	}

	return w
}

// Add adds a new keypair to the wallet and saves it
func (wallet *Wallet) Add(pubKey cryptography.PublicKey, secretKey *big.Int) {
	addr := Address{PublicKey: pubKey, SecretKey: secretKey}
	wallet.Addresses = append(wallet.Addresses, addr)
	wallet.Save()
}

// GenerateNew creates a new keypair and saves it
func (wallet *Wallet) GenerateNew() string {
	secretKey, pubKey := cryptography.RandomKeyPair()
	wallet.Add(pubKey, secretKey)
	return pubKey.ToAddress()
}

// ListAddresses returns a slice of the addresses in the wallet
func (wallet *Wallet) ListAddresses() (addresses []string) {
	addresses = make([]string, 0)
	for _, addr := range wallet.Addresses {
		addresses = append(addresses, addr.PublicKey.ToAddress())
	}
	return
}


func getWalletData() ([]byte, error) {
	if _, err := os.Stat(getWalletFolder()); os.IsNotExist(err) {
		os.Mkdir(getWalletFolder(), 0777)
	}
	if _, err := os.Stat(getWalletFile()); os.IsNotExist(err) {
		return nil, err
	}
	return ioutil.ReadFile(getWalletFile())
}

func getWalletFolder() string {
	cwd, _ := os.Getwd()
	return cwd + "/configs"
}

func getWalletFile() string {
	return getWalletFolder() + "/wallet.json"
}