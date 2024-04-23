package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	pk, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("private key: ", pk)

	addr := crypto.PubkeyToAddress(pk.PublicKey).String()
	fmt.Println(addr)

	return nil
}
