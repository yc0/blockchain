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
func (cli *CLI) Run() {
	cli.validateArgs()

	// addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	// addBlockData := addBlockCmd.String("data", "", "Block data")
	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockChainAddr := createBlockChainCmd.String("address", "", "The address to send genesis block reward to")
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
}
