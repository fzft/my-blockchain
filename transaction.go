package main

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"strconv"
	"time"
)

type Transaction struct {
	Hash        string
	fromAddress string  // the address of the sender
	toAddress   string  // the address of the recipient
	Amount      float64 // the amount being transferred
	Timestamp   int64   // the time the transaction was created (in Unix time)
	Signature   string  // the digital signature of the transaction
}

func NewTransaction(fromAddress string, toAddress string, amount float64) *Transaction {
	tx := &Transaction{
		fromAddress: fromAddress,
		toAddress:   toAddress,
		Amount:      amount,
		Timestamp:   time.Now().Unix(),
	}

	return tx
}

// getTransactionByteData returns the byte data of the transaction
func (t *Transaction) getTransactionByteData() []byte {
	fromAddress := t.fromAddress
	if t.fromAddress == "" {
		fromAddress = "0"
	}

	return []byte(fromAddress + t.toAddress + strconv.FormatFloat(t.Amount, 'f', -1, 64) + strconv.FormatInt(t.Timestamp, 10))
}

// calculateHash calculates the hash of the transaction
func (t *Transaction) calculateHash() common.Hash {
	hash := crypto.Keccak256Hash(t.getTransactionByteData())
	return hash
}

// signTransaction signs the transaction with the private key
func (t *Transaction) signTransaction(privateKeyHex string) error {

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Printf("Error HexToECDSA private key: %v", err)
		return err
	}

	hashTx := t.calculateHash()
	// Sign the hashed message using the private key
	signature, err := crypto.Sign(hashTx.Bytes(), privateKey)
	if err != nil {
		return err
	}

	t.Signature = hexutil.Encode(signature)

	log.Printf("Transaction signed successfully")
	return err
}

// isValid checks if the signature is valid
func (t *Transaction) isValid() bool {
	fromAddress := t.fromAddress
	if t.fromAddress == "" {
		fromAddress = "0"
		return true
	}

	signature, err := hexutil.Decode(t.Signature)
	if err != nil {
		log.Printf("Error decoding signature: %v", err)
		return false
	}

	hashTx := t.calculateHash()
	recoveredPublicKeyBytes, err := crypto.Ecrecover(hashTx[:], signature)

	// Create an ECDSA public key from the raw bytes using btcec library
	ecdsaPublicKey, _ := btcec.ParsePubKey(recoveredPublicKeyBytes)

	// Compress the ECDSA public key to its 33-byte format
	compressedPubKey := ecdsaPublicKey.SerializeCompressed()

	pubKey, err := crypto.DecompressPubkey(compressedPubKey)
	if err != nil {
		log.Printf("DecompressPubkey err:%v\n", err)
		return false
	}

	recoveredAddress := crypto.PubkeyToAddress(*pubKey).String()

	if fromAddress == recoveredAddress {
	} else {
		log.Println("Failed to recover public key")
		return false
	}

	signatureNoRecoverID := signature[:len(signature)-1] // remove recovery id
	valid := crypto.VerifySignature(recoveredPublicKeyBytes, hashTx.Bytes(), signatureNoRecoverID)
	return valid
}
