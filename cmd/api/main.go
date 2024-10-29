package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Index        int
	Timestamp    string
	Data         string
	PreviousHash []byte
	Hash         []byte
	Nonce        int
}

// CalculateHash computes the SHA-256 hash of the block
func (b *Block) CalculateHash() []byte {
	var hashData bytes.Buffer
	hashData.WriteString(strconv.Itoa(b.Index))
	hashData.WriteString(b.Timestamp)
	hashData.WriteString(b.Data)
	hashData.Write(b.PreviousHash)
	hashData.WriteString(strconv.Itoa(b.Nonce))

	hash := sha256.Sum256(hashData.Bytes())
	return hash[:]
}

// MineBlock performs proof-of-work
func (b *Block) MineBlock(difficulty int) {
	target := bytes.Repeat([]byte{0}, difficulty)
	for {
		b.Hash = b.CalculateHash()
		if bytes.HasPrefix(b.Hash, target) {
			fmt.Printf("Block mined: %x\n", b.Hash)
			break
		} else {
			b.Nonce++
		}
	}
}

// Blockchain is a series of validated Blocks
type Blockchain struct {
	Blocks     []*Block
	Difficulty int
}

// CreateGenesisBlock generates the first block
func CreateGenesisBlock() *Block {
	genesisBlock := &Block{
		Index:        0,
		Timestamp:    time.Now().String(),
		Data:         "Genesis Block",
		PreviousHash: []byte{},
		Nonce:        0,
	}
	genesisBlock.Hash = genesisBlock.CalculateHash()
	return genesisBlock
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := &Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().String(),
		Data:         data,
		PreviousHash: prevBlock.Hash,
		Nonce:        0,
	}
	newBlock.MineBlock(bc.Difficulty)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func main() {
	// Initialize blockchain with genesis block
	blockchain := &Blockchain{
		Blocks:     []*Block{CreateGenesisBlock()},
		Difficulty: 3, // Adjust the difficulty as needed
	}

	// Add new blocks
	blockchain.AddBlock("First Block after Genesis")
	blockchain.AddBlock("Second Block after Genesis")

	// Print the blockchain
	for _, block := range blockchain.Blocks {
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Timestamp: %s\n", block.Timestamp)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("PreviousHash: %x\n", block.PreviousHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Println("-------------------------------")
	}
}
