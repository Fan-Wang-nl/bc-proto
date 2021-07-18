package main

import (
	"flag"
	"fmt"
	"os"
)

type CLI struct {
	//bc *BlockChain
}

const usage = `
	createChain --address ADDRESS	"create a blockchain"
	addBlock 	--data Data	"add a block to this blockchain"
	send --from FROM --to TO --amount AMOUNT	"send coin from FROM to TO"
	getBalance	--address ADDRESS	"check balance of an address"
	printChain			"print all blocks"`

const AddBlockCmdString = "addBlock"
const PrintChainCmdString = "printChain"
const CreateChainCmdString = "createChain"

func (cli *CLI)parameterCheck() {
	if len(os.Args) < 2{
		cli.printUsage()
	}
}

func (cli *CLI)Run() {
	cli.parameterCheck()

	createChainCmd := flag.NewFlagSet(CreateChainCmdString, flag.ExitOnError)
	addBlockCmd := flag.NewFlagSet(AddBlockCmdString, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(PrintChainCmdString, flag.ExitOnError)

	addBlockPara := addBlockCmd.String("data", "","block transaction info")
	createChainPara := createChainCmd.String("address", "","address info")

	switch os.Args[1] {
	case CreateChainCmdString:
		err := createChainCmd.Parse(os.Args[2:])
		CheckErr("parse Run parameter 0", err)
		//after parsed, the data will be injected to addBlockPara
		if createChainCmd.Parsed(){
			if *createChainPara == ""{
				println("address should not be empty")
				cli.printUsage()
			}
			cli.CreateChain(*createChainPara)
		}
	case AddBlockCmdString:
		err := addBlockCmd.Parse(os.Args[2:])
		CheckErr("parse Run parameter 1", err)
		//after parsed, the data will be injected to addBlockPara
		if addBlockCmd.Parsed(){
			if *addBlockPara == ""{
				println("transaction data should not be empty")
				cli.printUsage()
			}
			cli.AddBlock(*addBlockPara)
		}
	case PrintChainCmdString:
		err := printChainCmd.Parse(os.Args[2:])
		CheckErr("parse Run parameter 2", err)
		if printChainCmd.Parsed(){
			cli.PrintChain()
		}
	default:
		println("unknown command")
		cli.printUsage()
	}
}

func (cli * CLI)printUsage() {
	fmt.Println(usage)
	os.Exit( 1)
}


