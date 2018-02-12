package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp    int64
	Transactions []*Transaction
	// Data         []byte
	Prev  []byte
	Hash  []byte
	Nonce int // counter
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

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
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
func NewBlock(transactions []*Transaction, prev []byte) *Block {
	// b := &Block{time.Now().Unix(), []byte(data), []byte{}, prev, 0}
	b := &Block{time.Now().Unix(), transactions, prev, []byte{}, 0}
	pow := NewProofOfWork(b)

	nonce, hash := pow.Run()
	b.Nonce = nonce
	b.Hash = hash[:]
	return b
}
func Genesis(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}
