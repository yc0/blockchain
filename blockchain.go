package main

import (
	"log"
	// As the result of bugs on golang/leveldb and instable,
	// I prefer to adopt syndtr/goleveldb one
	// "github.com/syndtr/goleveldb/leveldb"
	""
)

const (
	dbFile      = "blockchain.db"
	blockBucket = "blocks"
)

type BlockChain struct {
	// Blocks []*Block
	Tip []byte
	Db  *leveldb.DB // who should close the DB ?
}
type Iterator interface {
	Iterator() *BlockchainIterator
}
type BlockchainIterator struct {
	CurrentHash []byte
	Db          *leveldb.DB
}

func (bc *BlockChain) AddBlock(data string) {
	// prev := bc.Blocks[len(bc.Blocks)-1]
	// bc.Blocks = append(bc.Blocks, NewBlock(data, prev.Hash))

	lastHash, err := bc.Db.Get([]byte("l"), nil)
	// I prefer that the method which open a connection is responsible for closing a connection
	if err != nil {
		log.Panic(err)
	}
	_new := NewBlock(data, lastHash)
	batch := new(leveldb.Batch)
	batch.Put(_new.Hash, _new.Serialize())
	batch.Put([]byte("l"), _new.Hash)
	err = bc.Db.Write(batch, nil)
	if err != nil {
		log.Panic(err)
	}
	bc.Tip = _new.Hash
}

func NewBlockChain() *BlockChain {
	var tip []byte
	db, err := leveldb.OpenFile(dbFile+"/"+blockBucket, &opt.Options{ErrorIfMissing: true})
	if err != nil {
		db, err := leveldb.OpenFile(dbFile+"/"+blockBucket, nil)
		gensis := Genesis()
		batch := new(leveldb.Batch)
		batch.Put(gensis.Hash, gensis.Serialize())
		batch.Put([]byte("l"), gensis.Hash)
		err = db.Write(batch, nil)
		if err != nil {
			log.Panic(err)
		}
		tip = gensis.Hash
	} else {
		tip, err = db.Get([]byte("l"), nil)
		if err != nil {
			log.Panic(err)
		}
	}
	return &BlockChain{tip, db}
}

// Inspect blockchain
func (bc *BlockChain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.Tip, bc.Db}
}

func (i *BlockchainIterator) Next() *Block {

	encodedBlock, err := i.Db.Get(i.CurrentHash, nil)
	if err != nil {
		log.Panic(err)
	}
	block := DeserializeBlock(encodedBlock)
	i.CurrentHash = block.Prev
	return block
}
