package main

import (
	"encoding/hex"
	"log"
	// As the result of bugs on golang/leveldb and instable,
	// I prefer to adopt syndtr/goleveldb one
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const (
	dbFile              = "blockchain.db"
	blockBucket         = "blocks"
	genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
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

func (bc *BlockChain) MineBlock(transactions []*Transaction) {
	lastHash, err := bc.Db.Get([]byte("l"), nil)
	// I prefer that the method which open a connection is responsible for closing a connection
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

	batch := new(leveldb.Batch)
	batch.Put(newBlock.Hash, newBlock.Serialize())
	batch.Put([]byte("l"), newBlock.Hash)
	err = bc.Db.Write(batch, nil)
	if err != nil {
		log.Panic(err)
	}
	bc.Tip = newBlock.Hash
}

// func (bc *BlockChain) AddBlock(data string) {
// 	// prev := bc.Blocks[len(bc.Blocks)-1]
// 	// bc.Blocks = append(bc.Blocks, NewBlock(data, prev.Hash))

// 	lastHash, err := bc.Db.Get([]byte("l"), nil)
// 	// I prefer that the method which open a connection is responsible for closing a connection
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	_new := NewBlock(data, lastHash)
// 	batch := new(leveldb.Batch)
// 	batch.Put(_new.Hash, _new.Serialize())
// 	batch.Put([]byte("l"), _new.Hash)
// 	err = bc.Db.Write(batch, nil)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	bc.Tip = _new.Hash
// }

func NewBlockChain(address string) *BlockChain {
	var tip []byte
	db, err := leveldb.OpenFile(dbFile+"/"+blockBucket, &opt.Options{ErrorIfMissing: true})
	if err != nil {
		db, err := leveldb.OpenFile(dbFile+"/"+blockBucket, nil)
		coinbaseTx := NewCoinBaseTx(address, genesisCoinbaseData)
		gensis := Genesis(coinbaseTx)
		batch := new(leveldb.Batch)
		batch.Put(gensis.Hash, gensis.Serialize())
		batch.Put([]byte("l"), gensis.Hash)
		err = db.Write(batch, nil)
		if err != nil {
			log.Panic(err)
		}
		tip = gensis.Hash
		return &BlockChain{tip, db}
	} else {
		tip, err = db.Get([]byte("l"), nil)
		if err != nil {
			log.Panic(err)
		}
		return &BlockChain{tip, db}
	}
}

// Inspect blockchain
func (bc *BlockChain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.Tip, bc.Db}
}
func (bc *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction
	spentTXOs := make(map[string][]int)

	bci := bc.Iterator()

	for {
		block := bci.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			for outidx, out := range tx.Vout {
				// was the output spent ?
				// Bruce Eckel wrote in "Thinking in Java" following idea:
				// "It’s important to remember that the only reason to
				// use labels in Java is when you have nested loops
				// and you want to break or continue through more than one nested level."
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outidx { //spent
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTxs = append(unspentTxs, *tx)
				}

			}

			if !tx.IsCoinBaseTx() {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}
		if len(block.Prev) == 0 {
			break
		}
	}
	return unspentTxs
}

func (bc *BlockChain) FindUTXO(address string) []TXOutput {
	var ret []TXOutput
	unspentTxs := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTxs {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				ret = append(ret, out)
			}
		}
	}
	return ret
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
