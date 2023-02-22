package main

import (
	"bytes"
	"github.com/ethereum/go-ethereum/crypto"
)

type MerkleTree struct {
	Root *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	// Create a Merkle node for each data element
	for _, datum := range data {
		nodes = append(nodes, MerkleNode{nil, nil, datum})
	}

	// If the number of nodes is odd, duplicate the last node
	if len(nodes)%2 != 0 {
		nodes = append(nodes, nodes[len(nodes)-1])
	}

	// Build the tree by repeatedly hashing pairs of nodes
	for i := 0; i < len(nodes)/2; i++ {
		left := nodes[i*2]
		right := nodes[i*2+1]
		hash := crypto.Keccak256Hash(append(left.Data, right.Data...))
		parent := MerkleNode{&left, &right, hash[:]}
		nodes = append(nodes, parent)
	}

	return &MerkleTree{&nodes[len(nodes)-1]}
}

func TraverseMerkleTree(tree *MerkleTree, targetData []byte) (*MerkleNode, bool) {
	// Start at the root node
	node := tree.Root

	// Traverse down the tree until we find a leaf node containing the target data
	for node.Left != nil {
		leftHash := node.Left.Data
		rightHash := node.Right.Data
		targetHash := crypto.Keccak256Hash(append(leftHash, rightHash...))

		if bytes.Equal(targetData, node.Left.Data) {
			return node.Left, true
		} else if bytes.Equal(targetData, node.Right.Data) {
			return node.Right, true
		} else if bytes.Equal(targetData, targetHash[:]) {
			return node, true
		}

		// If the target data is in the left subtree, traverse down the left subtree
		if bytes.Compare(targetData, node.Left.Data) < 0 {
			node = node.Left
		} else {
			// Otherwise, traverse down the right subtree
			node = node.Right
		}
	}

	// If we didn't find the target data, return nil and false
	return nil, false
}

func GetMerkleRoot(tree *MerkleTree) []byte {
	return tree.Root.Data
}
