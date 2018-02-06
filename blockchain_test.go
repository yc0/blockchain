package main

import (
	"fmt"
	"strconv"
	"testing"
)

func TestBlockChainGenesis(t *testing.T) {
	bc := NewBlockChain()
	bc.AddBlock("Send 1 BTC to NSL")
	bc.AddBlock("Send 2 more BTC to NSL")
	for _, block := range bc.Blocks {
		fmt.Printf("Prev.Hash %x\n", block.Prev)
		fmt.Printf("Data %s\n", block.Data)
		fmt.Printf("Data %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW :%s\n", strconv.FormatBool(pow.Validate()))
		fmt.Print("\n\n")
	}
}
