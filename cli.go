package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {
	Bc *BlockChain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -address ADDRESS :Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  printchain :print all the blocks of the blockchain")
	fmt.Println("  getbalance -address ADDRESS :Get balance of ADDRESS")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT :Send AMOUNT of coins from FROM address to TO")
}
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}
func (cli *CLI) createBlockChain(address string) {
	bc := NewBlockChain(address)
	defer bc.Db.Close()
	fmt.Println("done !")
}

// func (cli *CLI) addBlock(data string) {
// 	cli.Bc.AddBlock(data)
// 	fmt.Println("Success !!")
// }
func (cli *CLI) printChain() {
	bc := NewBlockChain("")
	defer bc.Db.Close()

	bci := bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("Prev Hash: %x\n", block.Prev)
		// fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))

		if len(block.Prev) == 0 {
			break
		}
	}
}
func (cli *CLI) getBalance(address string) {
	bc := NewBlockChain(address)
	defer bc.Db.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

// Sending coins means creating a transaction and adding it to the blockchain via mining a block.
// But Bitcoin doesn’t do this immediately (as we do).
// Instead, it puts all new transactions into memory pool (or mempool),
// and when a miner is ready to mine a block,
// it takes all transactions from the mempool and creates a candidate block.
// Transactions become confirmed only when a block containing them is mined and added to the blockchain.
func (cli *CLI) send(from, to string, amount int) {
	bc := NewBlockChain(from)
	defer bc.Db.Close()

	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("Success!")
}
func (cli *CLI) Run() {
	cli.validateArgs()

	// addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	// addBlockData := addBlockCmd.String("data", "", "Block data")
	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockChainAddr := createBlockChainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	// case "addblock":
	// 	err := addBlockCmd.Parse(os.Args[2:])
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	// if addBlockCmd.Parsed() {
	// 	if *addBlockData == "" {
	// 		addBlockCmd.Usage()
	// 		os.Exit(1)
	// 	}
	// 	cli.addBlock(*addBlockData)
	// }

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockChainAddr == "" {
			createBlockChainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockChain(*createBlockChainAddr)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}
