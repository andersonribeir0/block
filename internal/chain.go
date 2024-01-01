package internal

import (
	"encoding/hex"
	"fmt"
	"time"
)

type Blockchain struct {
	blocks     []Block
	difficulty int
}

func NewBlockchain(difficulty int) *Blockchain {
	genesisBlock := Block{
		data:      []byte("Genesis Block"),
		timestamp: time.Now(),
	}

	genesisBlock.hash, _ = genesisBlock.CalculateHash()

	return &Blockchain{blocks: []Block{genesisBlock}, difficulty: difficulty}
}

func (b *Blockchain) IsValidBlockchain() bool {
	for i := 1; i < len(b.blocks); i++ {
		if !b.blocks[i].IsValidHash(b.difficulty) {
			return false
		}
	}

	return true
}

func (b *Blockchain) AddBlock(data []byte, kg *KeyGen) error {
	block, err := NewBlock(data, kg)
	if err != nil {
		return err
	}

	block.prevHash = b.blocks[len(b.blocks)-1].hash

	hash, err := block.CalculateHash()
	if err != nil {
		return err
	}

	block.hash = hash

	err = block.MineBlock(b.difficulty)
	if err != nil {
		return err
	}

	if !block.IsValid(b.difficulty, kg) {
		return fmt.Errorf("invalid signature")
	}

	b.blocks = append(b.blocks, *block)
	return nil
}

func (b *Blockchain) Log() {
	for i, block := range b.blocks {
		fmt.Printf("Block %d: PrevHash: %+v\n", i, hex.EncodeToString(block.prevHash))
		fmt.Printf("Block %d: Hash: %s\n", i, hex.EncodeToString(block.hash))
		fmt.Printf("Block %d: Data: %s\n", i, string(block.data))
		fmt.Printf("Block %d: Nonce: %d\n", i, block.nonce)
		fmt.Printf("Block %d: Signature: %s\n", i, hex.EncodeToString(block.signature))
		fmt.Printf("Block %d: Timestamp: %s\n", i, block.timestamp.String())
		fmt.Println("--------------------------------")
	}
}
