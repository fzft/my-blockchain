package main

import (
	"fmt"
	"log"
	"time"
)

type Blockchain struct {
	chain               []*Block
	difficulty          int
	pendingTransactions []*Transaction
	miningReward        int
}

func NewBlockchain() *Blockchain {
	blockChain := &Blockchain{
		difficulty:   2,
		miningReward: 100,
	}
	blockChain.chain = append(blockChain.chain, blockChain.createGenesisBlock())
	return blockChain
}

// createGenesisBlock creates the first block in the chain
func (bc *Blockchain) createGenesisBlock() *Block {

	dateString := "2023-01-01"
	date, _ := time.Parse("2006-01-02", dateString)

	genesisBlock := NewBlock(date.Unix(), []*Transaction{}, "0")
	genesisBlock.hash = "0"
	return genesisBlock
}

// minePendingTransactions mines the pending transactions and adds a new block to the chain
func (bc *Blockchain) minePendingTransactions(miningRewardAddress string) {
	rewardTx := NewTransaction("", miningRewardAddress, float64(bc.miningReward))
	bc.pendingTransactions = append(bc.pendingTransactions, rewardTx)
	block := NewBlock(time.Now().Unix(), bc.pendingTransactions, bc.getLatestBlock().hash)

	block.mineBlock(bc.difficulty)

	log.Println("Block mined successfully")
	bc.chain = append(bc.chain, block)
	bc.pendingTransactions = []*Transaction{}
}

// addTransaction adds a new transaction to the pending transactions
func (bc *Blockchain) addTransaction(transaction *Transaction) (err error) {
	if transaction.fromAddress == "" || transaction.toAddress == "" {
		err = fmt.Errorf("transaction must include a sender and a recipient")
		log.Fatalf("%v", err)
		return
	}

	if !transaction.isValid() {
		err = fmt.Errorf("transaction is not valid")
		log.Fatalf("%v", err)
		return
	}

	if transaction.Amount <= 0 {
		err = fmt.Errorf("transaction amount must be greater than 0")
		log.Fatalf("%v", err)
		return
	}

	walletBalance := bc.getBalanceOfAddress(transaction.fromAddress)
	if walletBalance < transaction.Amount {
		err = fmt.Errorf("not enough balance")
		return
	}

	pendingTxForWallet := filter(bc.pendingTransactions, func(tx *Transaction) bool {
		return tx.fromAddress == transaction.fromAddress
	})

	if len(pendingTxForWallet) > 0 {
		totalPendingAmount := reducer(pendingTxForWallet)
		totalAmount := totalPendingAmount + transaction.Amount
		if totalAmount > walletBalance {
			err = fmt.Errorf("pending transactions for this wallet is higher than its balance")
			log.Fatalf("%v", err)
			return
		}
	}

	bc.pendingTransactions = append(bc.pendingTransactions, transaction)
	log.Printf("transaction added %+v", transaction)
	return
}

func (bc *Blockchain) getBalanceOfAddress(address string) float64 {
	var balance float64
	for _, block := range bc.chain {

		for _, trans := range block.transactions {
			if trans.fromAddress == address {
				balance -= trans.Amount
			}

			if trans.toAddress == address {
				balance += trans.Amount
			}
		}
	}

	return balance
}

// getAllTransactionsForWallet returns all transactions for a given wallet address
func (bc *Blockchain) getAllTransactionsForWallet(address string) []*Transaction {
	var txs []*Transaction
	for _, block := range bc.chain {
		for _, tx := range block.transactions {
			if tx.fromAddress == address || tx.toAddress == address {
				txs = append(txs, tx)
			}
		}
	}

	return txs
}

// isChainValid checks if the chain is valid
func (bc *Blockchain) isChainValid() bool {

	// check if the genesis block is valid
	realGenesis := bc.createGenesisBlock()

	if realGenesis.hash != bc.chain[0].hash {
		return false
	}

	for i := 1; i < len(bc.chain); i++ {
		currentBlock := bc.chain[i]
		previousBlock := bc.chain[i-1]

		if currentBlock.hash != currentBlock.calculateHash() {
			return false
		}

		if currentBlock.previousHash != previousBlock.hash {
			return false
		}

		if !currentBlock.hasValidTransactions() {
			return false
		}
	}

	return true
}

// getLatestBlock returns the latest block in the chain
func (bc *Blockchain) getLatestBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) queryTransactionInBlock(hash string) *Transaction {

	return nil
}