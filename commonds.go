package main

import "fmt"

func (cli *CLI)AddBlock(data string) {
	bc := getBlockChainHandler()
	bc.AddBlock(data)
}


func (cli *CLI)PrintChain() {
	bc := getBlockChainHandler()
	it := bc.NewIterator()
	for{
		block := it.Next()
		fmt.Printf("Current Has: %x\n", block.Hash)
		fmt.Printf("Previous Has: %x\n", block.PreviousBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Is Valid: %v\n", NewProofOfWork(block).IsValid())

		if string(block.Data) == "Genesis Block"{
			println("over!")
			break
		}
	}
}


func (cli *CLI) CreateChain(address string) {
	bc := InitBlockChain()
	bc.db.Close()
	println("create a blockchain successfully")


}