package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
)

type Block struct {
	timestamp int64
	transactions []*Transaction
	previousHash string
	hash string
	nonce int
}

// NewBlock creates and returns a new block
func NewBlock(timestamp int64, transactions []*Transaction, previousHash string) *Block {
	block := &Block{
		timestamp: timestamp,
		transactions: transactions,
		previousHash: previousHash,
		nonce: 0,
	}
	block.hash = block.calculateHash()
	return block
}

func (b *Block) getBlockByteData() []byte {
	txBytes , _ := json.Marshal(b.transactions)
	return []byte(  string(b.previousHash) + strconv.Itoa(int(b.timestamp) ) + string(txBytes)   + strconv.Itoa(b.nonce))
}

func (b *Block) calculateHash() string {
	hash := sha256.Sum256(b.getBlockByteData())
	return hex.EncodeToString(hash[:])
}

func (b *Block) mineBlock(difficulty int) {
	prefix := ""
	for i := 0; i < difficulty; i++ {
		prefix += "0"
	}
	for !strings.HasPrefix(b.hash, prefix) {
		b.nonce++
		b.hash = b.calculateHash()
	}
}

func (b *Block) hasValidTransactions() bool {
	for _, tx := range b.transactions {
		if !tx.isValid() {
			return false
		}
	}
	return true
}


