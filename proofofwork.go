package blockchain

import (
	"crypto/sha256"
	"fmt"
	"math"
	"bytes"
	"math/big"
)

const (
	targetBits = 16
	maxNonce = math.MaxInt64
)

type Proof_Of_Work struct {
	block *Block
	target *big.Int // a struct type
}

func (pow *Proof_Of_Work) New(b *Block) *Proof_Of_Work {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits)) // sha 256 left shift
	pow = &Proof_Of_Work{b,target}
	return pow
}

func (pow *Proof_Of_Work) prepareData(nonce int) []byte {
	/*
	 * :type nounce: int counter 
	 */
	data := bytes.Join(
		[][]byte{
			pow.block.Prev,
			pow.block.Data,
			Int2Hex(pow.block.Timestamp),
			Int2Hex(int64(targetBits)),
			Int2Hex(int64(nonce)),
		}, []byte{},
	)
	return data
}


/**
 run performs a proof-of-work
 **/
func (pow *Proof_Of_Work) Run() (int,[]byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	fmt.Printf("Mining the block containing '%s' \n",pow.block.Data)
	defer fmt.Printf("\n\n")
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 { // if hashInt < target -->valid	
			break
		} else {
			nonce++
		}
	}
	return nonce, hash[:]
}

/**
 * validate block's PoW
 */
func (pow *Proof_Of_Work) Validate() bool {
	var hashInt big.Int
	// fmt.Println(".....",pow.block.Nonce)
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}