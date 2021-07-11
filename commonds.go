package main

import "fmt"

func (cli *CLI)AddBlock(data string) {
	cli.bc.AddBlock(data)
}

func (cli *CLI)PrintChain() {
	it := cli.bc.NewIterator()
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