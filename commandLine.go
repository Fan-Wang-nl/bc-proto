package main

import (
	"flag"
	"fmt"
	"os"
)

type CLI struct {
	bc *BlockChain
}

const usage = `
	addBlock --data Data "add a block to this blockchain"
	printChain			 "print all blocks"`

const AddBlockCmdString = "addBlock"
const PrintChainCmdString = "printChain"

func (cli *CLI)parameterCheck() {
	if len(os.Args) < 2{
		cli.printUsage()
	}
}

func (cli *CLI)Run() {
	cli.parameterCheck()
	addBlockCmd := flag.NewFlagSet(AddBlockCmdString, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(PrintChainCmdString, flag.ExitOnError)

	addBlockPara := addBlockCmd.String("data", "","block transaction info")

	switch os.Args[1] {
	case AddBlockCmdString:
		err := addBlockCmd.Parse(os.Args[2:])
		CheckErr("parse Run parameter 1", err)
		//after parsed, the data will be injected to addBlockPara
		if addBlockCmd.Parsed(){
			if *addBlockPara == ""{
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
		cli.printUsage()

	}
}

func (cli * CLI)printUsage() {
	fmt.Println(usage)
	os.Exit( 1)
}
