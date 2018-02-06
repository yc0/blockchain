package main

func main() {
	bc := NewBlockChain()
	defer bc.Db.Close()
	cli := CLI{bc}
	// cli.printUsage()
	cli.Run()
}
