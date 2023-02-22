package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
)

func GenerateKeyPair() (privateKeyHex string, publicKeyHex string, address string, err error) {
	// Generate a private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
		return
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex = hex.EncodeToString(privateKeyBytes)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
		return
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyHex = hex.EncodeToString(publicKeyBytes)

	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return
}
