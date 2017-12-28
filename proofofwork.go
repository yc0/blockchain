package blockchain

import (
	"sync"
	"crypto/sha256"
	"fmt"
	"bytes"
	"math/big"
	"runtime"
)

const (
	targetBits = 16
	maxNonce = 1<<24 - 1
	maxConcurrencies = 32
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
	runtime.GOMAXPROCS(maxConcurrencies)
	var hashInt big.Int
	var locker sync.Once
	isDone := false // a sort of optimisitic lock
	c_n := make(chan int)
	// done := make(chan bool, maxConcurrencies-1)
	fmt.Printf("Mining the block containing '%s' \n",pow.block.Data)
	defer fmt.Printf("\n\n")

	step := (1<<16)
	size := maxNonce/step
	mod := maxNonce%step

	concurrentGoroutines := make(chan struct{}, maxConcurrencies)
	for i := 0; i < maxConcurrencies; i++ {
		concurrentGoroutines <- struct {}{}
	}
	for i:= 0; i <=size; i+=1 {
		min := i*step
		var max int
		if min+mod == maxNonce {
			max = maxNonce
		} else {
			max = step*(i+1)-1
		}
		go func(min, max int, cn chan int) {
			<-concurrentGoroutines
			for i:=min; i <= max && !isDone; i++ {
				data := pow.prepareData(i)
				h  := sha256.Sum256(data)
				// fmt.Printf("\r%x", hash)
				hashInt.SetBytes(h[:])
				if hashInt.Cmp(pow.target) == -1 {
					locker.Do(func() {
						cn <- i
						defer close(cn)
					})
					// for {
					// 	select {
					// 	case cn <- i:
					// 		defer close(cn)
					// 		return
					// 	case <- done:
					// 		return
					// 	}
					// }
				}
			}
			// done <- true
		}(min, max, c_n)
		concurrentGoroutines <- struct{}{}
	}
	fmt.Println("awaiting..")
	nonce := <-c_n
	// for i := 0; i < maxConcurrencies-1; i++ {
	// 	done <- true
	// }
	// defer close(done)
	isDone = true
	fmt.Printf("%d ::", nonce)
	data := pow.prepareData(nonce)
	hash := sha256.Sum256(data)
	fmt.Printf("%x \n", hash)
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