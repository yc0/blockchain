# BlockChain 
Building blockchain in Golang

refer to [jeiwan blockchain building tutorial]( https://jeiwan.cc/posts/building-blockchain-in-go-part-1/)
## Part-1 Branch
build the simple block with sha256 hash; however, I ignore this part, and move toward part-2 directly

## Part-2 Branch
build the proof of work hash for block.
The different part is that I utilize the goroutine, sync.once and limit goroutine to expedite the pow calculation.

This modification gave me runable and reasonable time-process for PoW on my laptop (new macbook)

Here are my constanst parameters

```
targetBits = 26 
maxNonce = 1<<28 - 1
maxConcurrencies = 32
```
### Workable Parameter 
While trying to test and implement your blockchain, you might need reasonable parameters, saying below 1.0s

Here are the suitable configurations.
```
targetBits = 16
maxNonce = 1<<24 - 1
maxConcurrencies = 32
```

## Part-3 Persistence
Bitcoin Core uses two “buckets” to store data:

1. blocks stores metadata describing all the blocks in a chain
2. chainstate stores the state of a chain, which is all currently unspent transaction outputs and some metadata.

Also, blocks are stored as separate files on the disk. This is done for a performance purpose: reading a single block won’t require loading all (or some) of them into memory. We won’t implement this.

In blocks, the key -> value pairs are:

> 'b' + 32-byte block hash -> block index record.
>
> 'f' + 4-byte file number -> file information record
>
> 'l' -> 4-byte file number: the last block file number used
>
> 'R' -> 1-byte boolean: whether we're in the process of reindexing
>
> 'F' + 1-byte flag name length + flag name string -> 1 byte boolean: 
> various flags that can be on or off
>
> 't' + 32-byte transaction hash -> transaction index record

In chainstate, the key -> value pairs are:

> 'c' + 32-byte transaction hash -> unspent transaction output record for that transaction
> 
> 'B' -> 32-byte block hash: the block hash up to which the database represents the unspent transaction outputs
>(Detailed explanation can be found [here](https://en.bitcoin.it/wiki/Bitcoin_Core_0.11_(ch_2):_Data_Storage))


Since we have **NO transactions** yet, we’re going to have only blocks bucket. Also, as said above, we will store the whole DB **as a single file**, without storing blocks in separate files. So we won’t need anything related to file numbers. So these are key -> value pairs we’ll use:

> 32-byte block-hash -> Block structure (serialized)
>
> 'l' -> the hash of the last block in a chain

That’s all we need to know to start implementing the persistence mechanism.

### Serialization

We adopt levelDB, as the result of bitcoin Core where it was initially published by Satoshi Nakamoto and is currently a reference implementation of Bitcoin, uses LevelDB (although it was introduced to the client only in 2012). 

In levelDB, the data can be only of []byte type, thus we'll use encoding/gob to serialize the structs

Meanwhile, official leveldb ported on golang repo says it is still instable, and most importantly, NOBODY contributes to the project for a year. Therefore, I adopt syndtr levelDB version.

In the implmentation, I do not implement transaction as Jeiwan did with boltDB and goroutine. Thus, I persist data through file system structure using folder as blockchain, and filename as blocks.

## Part-4 Transaction

Transactions are the heart of de-centralization and the only purpose of blockchain is to store transactions in a secure and reliable way, so no one could modify them after they are created. 

### General Mechanism

In de-centralization transactions like Bitcoin, payments are realized in completely different way. There are:

> 1. No accounts.
> 2. No balances.
> 3. No addresses.
> 4. No coins.
> 5. No senders and receivers.

Note that:
- some outputs are not linked to inputs.
- In **one transaction**, inputs can reference outputs from multiple transactions.
- An input must reference an output.

