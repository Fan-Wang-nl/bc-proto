package main

import "fmt"

func main(){
	bc := NewBlockChain()
	bc.AddBlock("A send to B 1 BTC")
	bc.AddBlock("B send to C 2 BTC")

	for _,block := range bc.blocks{
		fmt.Printf("Previous Has: %x\n", block.PreviousBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Is Valid: %v\n", NewProofOfWork(block).IsValid())
	}
}