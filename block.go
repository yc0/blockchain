package blockchain

import (
	"time"
)

type Block struct {
	Timestamp int64
	Data      []byte
	Hash     	[]byte
	Prev 			[]byte
	Nonce     int   // counter
	// merkle 		Merkle
}

// func New(data string) *Block {
// 	block := &Block {time.Now().Unix(), []byte(data), nil, nil}
// 	return block
// }
func (b *Block) New(data string, prevBlockHash []byte) *Block {
	b = &Block{time.Now().Unix(),[]byte(data),[]byte{}, prevBlockHash,0}
	pow := (&Proof_Of_Work{}).New(b)
	// pow.New(b)

	nonce, hash := pow.Run()
	b.Nonce = nonce
	b.Hash = hash[:]
	return b
}

func Genesis() *Block {
	var b *Block
	return b.New("Genesis Block", []byte{})
}