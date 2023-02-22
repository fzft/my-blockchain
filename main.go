package main

import (
	"log"
)

func main() {
	privateKey, pubKey, addr, err := GenerateKeyPair()
	if err != nil {
		return
	}

	log.Printf("privateKey: %v, pubKey: %v, addr: %v", privateKey, pubKey, addr)

	_, _, addr1, err := GenerateKeyPair()
	if err != nil {
		return
	}

	log.Println("addr1: ", addr1)
	_, _, addr2, err := GenerateKeyPair()

	if err != nil {
		return
	}
	log.Println("addr2: ", addr2)

	// create a new instance of blockchain
	blockchain := NewBlockchain()

	blockchain.minePendingTransactions(addr)

	log.Println("create a transaction & sign it with the private key")

	// create a transaction & sign it with the private key
	tx1 := NewTransaction(addr, addr1,100)

	// sign the transaction with the private key
	err = tx1.signTransaction(privateKey)
	if err != nil {
		log.Printf("Error signing transaction: %v", err)
		return
	}

	err = blockchain.addTransaction(tx1)
	if err != nil {
		log.Fatalf("Error adding transaction1: %v", err)
	}

	blockchain.minePendingTransactions(addr)

	log.Println("create another transaction & sign it with the private key")

	// create another transaction & sign it with the private key
	tx2 := NewTransaction(addr, addr2, 50)
	tx2.signTransaction(privateKey)
	err = blockchain.addTransaction(tx2)
	if err != nil {
		log.Fatalf("Error adding transaction2: %v", err)
	}

	blockchain.minePendingTransactions(addr)
	balance := blockchain.getBalanceOfAddress(addr)
	log.Println("balance: ", balance)

	isChainValid := blockchain.isChainValid()
	log.Println("isChainValid: ", isChainValid)


}
