package blockchain

import (
	"fmt"
	"testing"
)

func TestPow(t *testing.T) {
	var pow *Proof_Of_Work
	pow = pow.New(nil)
	fmt.Printf("%x\n",pow.target)
	
}