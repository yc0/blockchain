package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
)

const (
	subsidy = 10
)

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

// Let’s sum it up. Outputs are where “coins” are stored.
// Each output comes with an unlocking script,
// which determines the logic of unlocking the output.
// Every new transaction must have at least one input and output.
// An input references an output from a previous transaction and provides data
// (the ScriptSig field) that is used in the output’s unlocking script
// to unlock it and use its value to create new outputs.

type TXOutput struct {
	Value        int
	ScriptPubKey string
}

func NewCoinBaseTx(to, data string) *Transaction {
	txin := TXInput{[]byte{}, -1, ""}
	txout := TXOutput{subsidy, to}
	tx := &Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()
	return tx
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx *Transaction) IsCoinBaseTx() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

func NewUTXOTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	// Before creating new outputs, we first have to find all unspent outputs
	// and ensure that they store enough value.
	// This is what FindSpendableOutputs
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("Error: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		_txid, _ := hex.DecodeString(txid)
		for _, out := range outs {
			input := TXInput{_txid, out, from}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	// One that’s locked with the receiver address.
	// This is the actual transferring of coins to other address
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		// One that’s locked with the sender address.
		// This is a change.
		// It’s only created when unspent outputs hold more value than required
		// for the new transaction.
		// Remember: outputs are indivisible.
		outputs = append(outputs, TXOutput{amount - acc, from}) // a change
	}
	tx := &Transaction{nil, inputs, outputs}
	tx.SetID()
	return tx
}
