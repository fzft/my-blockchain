package main

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"strings"
	"time"
)

type Block struct {
	timestamp int64
	transactions []*Transaction
	previousHash string
	hash string
	nonce int

	merKleRoot string
}

// NewBlock creates and returns a new block
func NewBlock(timestamp int64, transactions []*Transaction, previousHash string) *Block {
	block := &Block{
		timestamp: timestamp,
		transactions: transactions,
		previousHash: previousHash,
		nonce: 0,
	}

	// exclude genesis block
	if previousHash != "0" {
		block.merKleRoot = block.getMerKleRoot()
		block.hash = block.calculateHash()
	}

	return block
}


// isTransactionIn checks if a transaction is in a block
func (b *Block) isTransactionIn(transaction *Transaction) bool {
	tree := b.buildMerKleTree()
	_, in := TraverseMerkleTree(tree, transaction.getTransactionByteData())
	return in
}


func (b *Block) buildMerKleTree() *MerkleTree {
	log.Printf("start to build tree")
	var txBytes [][]byte
	for _, tx := range b.transactions {
		txBytes = append(txBytes, tx.getTransactionByteData())
	}
	tree := NewMerkleTree(txBytes)
	log.Println("build tree success")
	return tree
}

// getMerKleRoot returns the merkle root of a block
func (b *Block) getMerKleRoot() string {
	tree := b.buildMerKleTree()
	rootData := GetMerkleRoot(tree)
	return hex.EncodeToString(rootData)
}

// getBlockByteData returns the byte data of a block, use merkle tree to optimize
func (b *Block) getBlockByteData() []byte {
	return []byte(string(b.nonce)  +  b.merKleRoot)
}

func (b *Block) calculateHash() string {
	hash := crypto.Keccak256Hash(b.getBlockByteData())
	return hex.EncodeToString(hash[:])[2:]
}

func (b *Block) mineBlock(difficulty int) {
	prefix := ""
	for i := 0; i < difficulty; i++ {
		prefix += "0"
	}
	for !strings.HasPrefix(b.hash, prefix) {
		b.nonce++
		b.hash = b.calculateHash()
		time.Sleep(time.Millisecond * 100)
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


