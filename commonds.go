package main

import "fmt"

func (cli *CLI)PrintChain() {
	bc := getBlockChainHandler()
	defer bc.db.Close()
	it := bc.NewIterator()
	for{
		block := it.Next()
		fmt.Printf("Current Hash: %x\n", block.Hash)
		fmt.Printf("Previous Hash: %x\n", block.PreviousBlockHash)
		fmt.Printf("transaction number: %d\n", len(block.Transactions))
		for _,transactionP := range (*block).Transactions{
			for _,input := range (*transactionP).TXInputs{
				fmt.Printf("	transaction input: %s;	%d; %s\n", string(input.TXID), input.Vout ,input.ScriptSig)
			}

			for _,output := range (*transactionP).TXOutputs{
				fmt.Printf("	transaction output: %f;	%s\n", output.Value ,output.ScriptPubKey)
			}
		}


		if len(block.PreviousBlockHash) == 0{
			println("over!")
			break
		}
	}
}

func (cli *CLI) GetBalance(address string){
	bc := getBlockChainHandler()
	defer bc.db.Close()
	utxos := bc.findUTXO(address)

	var total float64 = 0
	for _,utxo := range utxos{
		total += utxo.Value
	}

	fmt.Printf("%s has %f BTCs now\n", address,total)
}

func (cli *CLI) CreateChain(address string) {
	bc := InitBlockChain(address)
	bc.db.Close()
	println("create a blockchain successfully")
}

func (cli *CLI)send(from, to string, amount float64){
	bc := getBlockChainHandler()
	defer bc.db.Close()

	tx := NewTransaction(from, to, amount, bc)

	bc.AddBlock([]*Transaction{tx})

	println("send successfully")
}