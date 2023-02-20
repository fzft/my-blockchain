package main

import "time"

type BlockChain struct {
	chain []*Block
	difficulty int
	pendingTransactions []*Transaction
	miningReward int
}

func NewBlockChain() *BlockChain {
	blockChain := &BlockChain{
		difficulty: 2,
		miningReward: 100,
	}
	blockChain.createGenesisBlock()
	return blockChain
}

// createGenesisBlock creates the first block in the chain
func (bc *BlockChain) createGenesisBlock() {

	dateString := "2023-01-01"
	date, _ := time.Parse("2006-01-02", dateString)

	genesisBlock := NewBlock(date.Unix(), []*Transaction{}, "0")
	genesisBlock.hash = "0"
	bc.chain = append(bc.chain, genesisBlock)
}

// minePendingTransactions mines the pending transactions and adds a new block to the chain
func (bc *BlockChain) minePendingTransactions(miningRewardAddress string) {
	rewardTx := NewTransaction(nil, miningRewardAddress, float64(bc.miningReward))
	bc.pendingTransactions = append(bc.pendingTransactions, rewardTx)

	block := NewBlock(time.Now().Unix(), bc.pendingTransactions, bc.chain[len(bc.chain)-1].hash)
	bc.chain = append(bc.chain, block)

	block.mineBlock(bc.difficulty)
}

