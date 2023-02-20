package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

type Transaction struct {
	Sender    *string  // the address of the sender
	Recipient string  // the address of the recipient
	Amount    float64 // the amount being transferred
	Timestamp int64   // the time the transaction was created (in Unix time)
	Signature string  // the digital signature of the transaction
}

func NewTransaction(sender *string, recipient string, amount float64) *Transaction {
	return &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
		Timestamp: time.Now().Unix(),
	}
}

// getTransactionByteData returns the byte data of the transaction
func (t *Transaction) getTransactionByteData() []byte {
	fromAddress := *t.Sender
	if t.Sender == nil {
		fromAddress = "0"
	}
	return []byte(fromAddress + t.Recipient + strconv.FormatFloat(t.Amount, 'f', -1, 64) + strconv.FormatInt(t.Timestamp, 10))
}

// calculateHash calculates the hash of the transaction
func (t *Transaction) calculateHash() []byte {
	hash := sha256.Sum256(t.getTransactionByteData())
	return hash[:]
}

// signTransaction signs the transaction with the private key
func (t *Transaction) signTransaction(privateKey ecdsa.PrivateKey) error {
	pubKey := privateKey.Public().(*ecdsa.PublicKey)
	if !pubKey.Equal(t.Sender) {
		return fmt.Errorf("you cannot sign transactions for other wallets")
	}

	hashTx := t.calculateHash()
	// Sign the hashed message using the private key
	signature, err := privateKey.Sign(rand.Reader, hashTx, nil)
	if err != nil {
		return err
	}

	// Encode the signature in DER format
	der, err := asn1.Marshal(signature)
	if err != nil {
		return err
	}

	t.Signature = string(der)
	return err
}

// isValid checks if the signature is valid
func (t *Transaction) isValid() bool {
	fromAddress := *t.Sender
	if t.Sender == nil {
		fromAddress = "0"
	}

	// Decode the signature from DER format
	var decodedSignature struct {
		R, S *big.Int
	}

	_, err := asn1.Unmarshal([]byte(t.Signature), &decodedSignature)
	if err != nil {
		return false
	}

	pubKey, err := GetPublicKeyFromString(fromAddress)
	if err != nil {
		return false
	}

	valid := ecdsa.Verify(pubKey, t.calculateHash(), decodedSignature.R, decodedSignature.S)
	return valid
}

func GetPublicKeyFromString(keyStr string) (*ecdsa.PublicKey, error) {
	// Decode the string into a byte array
	keyBytes, _ := pem.Decode([]byte(keyStr))
	if keyBytes == nil {
		return nil, errors.New("invalid key string")
	}

	// Parse the byte array into a certificate
	cert, err := x509.ParseCertificate(keyBytes.Bytes)
	if err != nil {
		return nil, err
	}

	// Extract the public key from the certificate
	pubKey, ok := cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to extract public key")
	}

	return pubKey, nil
}

