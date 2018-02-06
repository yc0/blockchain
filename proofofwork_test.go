package main

import (
	"fmt"
	"testing"
)

func TestPow(t *testing.T) {
	pow := NewProofOfWork(nil)
	fmt.Printf("%x\n", pow.target)
}
