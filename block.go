package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp int64
	Data      []byte
	Hash      []byte
	Prev      []byte
	Nonce     int // counter
	// merkle 		Merkle
}

// func New(data string) *Block {
// 	block := &Block {time.Now().Unix(), []byte(data), nil, nil}
// 	return block
// }
func (b *Block) Serialize() []byte {
	var rst bytes.Buffer
	encoder := gob.NewEncoder(&rst)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return rst.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var b Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&b)
	if err != nil {
		log.Panic(err)
	}
	return &b
}

// follow up the convention form regarding New Object instead of New method
// func (b *Block) New(data string, prevBlockHash []byte) *Block {
// 	b = &Block{time.Now().Unix(), []byte(data), []byte{}, prevBlockHash, 0}
// 	pow := (&Proof_Of_Work{}).New(b)
// 	// pow.New(b)

// 	nonce, hash := pow.Run()
// 	b.Nonce = nonce
// 	b.Hash = hash[:]
// 	return b
func NewBlock(data string, prev []byte) *Block {
	b := &Block{time.Now().Unix(), []byte(data), []byte{}, prev, 0}
	pow := NewProofOfWork(b)

	nonce, hash := pow.Run()
	b.Nonce = nonce
	b.Hash = hash[:]
	return b
}
func Genesis() *Block {
	return NewBlock("Genesis Block", []byte{})
}
