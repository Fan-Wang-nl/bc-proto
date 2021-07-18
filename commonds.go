package main

import "fmt"

func (cli *CLI)AddBlock(data string) {
	//bc := getBlockChainHandler()
	//bc.AddBlock(data) //TODO
}


func (cli *CLI)PrintChain() {
	bc := getBlockChainHandler()
	it := bc.NewIterator()
	for{
		block := it.Next()
		fmt.Printf("Current Hash: %x\n", block.Hash)
		fmt.Printf("Previous Hash: %x\n", block.PreviousBlockHash)
		fmt.Printf("transaction number: %d\n", len(block.Transactions))
		fmt.Printf("Is Valid: %v\n", NewProofOfWork(block).IsValid())

		if len(block.PreviousBlockHash) == 0{
			println("over!")
			break
		}
	}
}


func (cli *CLI) CreateChain(address string) {
	bc := InitBlockChain(address)
	bc.db.Close()
	println("create a blockchain successfully")


}