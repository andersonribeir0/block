package main

import (
	"fmt"
	"time"

	"github.com/andersonribeir0/blocker/internal"
)

func main() {
	blockchain := internal.NewBlockchain(3)
	kg := internal.NewKeyGen()
	for i := 0; i < 100; i++ {
		block := fmt.Sprintf("block%d", i)
		fmt.Println("adding ", block)
		now := time.Now()
		blockchain.AddBlock([]byte(block), kg)
		elapsed := time.Since(now)
		fmt.Printf("time elapsed for %s: %f\n", block, elapsed.Seconds())
	}

	blockchain.Log()
}
