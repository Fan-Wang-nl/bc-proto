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
	send --from FROM --to TO --amount AMOUNT	"send coin from FROM to TO"
	getBalance	--address ADDRESS	"check balance of an address"
	printChain			"print all blocks"`

const PrintChainCmdString = "printChain"
const CreateChainCmdString = "createChain"
const GetBalanceCmdString = "getBalance"
const sendCmdString = "send"

func (cli *CLI)parameterCheck() {
	if len(os.Args) < 2{
		cli.printUsage()
	}
}

func (cli *CLI)Run() {
	cli.parameterCheck()

	createChainCmd := flag.NewFlagSet(CreateChainCmdString, flag.ExitOnError)
	//addBlockCmd := flag.NewFlagSet(AddBlockCmdString, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(PrintChainCmdString, flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet(GetBalanceCmdString, flag.ExitOnError)
	sendCmd := flag.NewFlagSet(sendCmdString, flag.ExitOnError)

	//addBlockPara := addBlockCmd.String("data", "","block transaction info")
	createChainPara := createChainCmd.String("address", "","address info")
	getBalancePara := getBalanceCmd.String("address", "","address info")

	//related to send
	sendParaFrom := sendCmd.String("from", "","sender address info")
	sendParaTo := sendCmd.String("to", "","receiver address info")
	sendParaAmount := sendCmd.Float64("amount", 0,"transaction amount")

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
	case sendCmdString:
		err := sendCmd.Parse(os.Args[2:])
		CheckErr("parse Run parameter 1", err)
		if *sendParaFrom == "" || *sendParaTo == "" || *sendParaAmount <= 0{
			println("invalid parameter")
			cli.printUsage()
		}
		cli.send(*sendParaFrom, *sendParaTo, *sendParaAmount)

	case PrintChainCmdString:
		err := printChainCmd.Parse(os.Args[2:])
		CheckErr("parse Run parameter 2", err)
		if printChainCmd.Parsed(){
			cli.PrintChain()
		}
	case GetBalanceCmdString:
		err := getBalanceCmd.Parse(os.Args[2:])
		CheckErr("parse Run parameter 3", err)
		//after parsed, the data will be injected to addBlockPara
		if getBalanceCmd.Parsed(){
			if *getBalancePara == ""{
				println("address should not be empty")
				cli.printUsage()
			}
			cli.GetBalance(*getBalancePara)
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


