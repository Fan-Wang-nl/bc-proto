package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"os"
)

const reward = 12.5
type Transaction struct {
	TXID      []byte
	TXInputs  []TXInput
	TXOutputs []TXOutput
}

type TXInput struct {
	//transaction ID of previous output
	TXID      []byte
	//index of the output, because a transaction may have multiple outputs
	Vout      int64
	//unlock script, specify who(private key) and how to use it
	ScriptSig string
}

type TXOutput struct {
	//payment value
	Value        float64
	//lock script, specify the address of the receiver
	ScriptPubKey string
}

func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction{
	//map[string][]int64; key is the transaction ID, and the array the indices inside the transaction
	validUTXOs := make(map[string][]int64)
	var total float64
	validUTXOs, total = bc.findSuitableUTXOs(from, amount)

	if total < amount{
		println("The UTXO is not enough")
		os.Exit(1)
	}

	var inputs []TXInput
	var outputs []TXOutput

	//1, create inputs. transform from outputs to inputs
	//traverse all transactions
	for TxID, outputIndices := range validUTXOs{
		//traverse all valid outputs of the transaction
		for _,index := range outputIndices{
			input := TXInput{TXID: []byte(TxID), Vout: index,  ScriptSig: from}
			inputs = append(inputs, input)
		}
	}

	//2 create outputs.
	output := TXOutput{amount, to}
	outputs = append(outputs, output)

	//return the change
	if total > amount{
		output := TXOutput{total - amount, from}
		outputs = append(outputs, output)
	}

	tx := Transaction{nil, inputs, outputs}
	tx.setTXID()
	return &tx
}

//check if current user owns the UTXO and can consume them
func (input *TXInput)validateUTXO(unlockScript string)  bool{
	return input.ScriptSig == unlockScript
}

//check if current user owns the UTXO
func (output *TXOutput)validateUTXO(unlockScript string)  bool{
	return output.ScriptPubKey == unlockScript
}

//set the transaction ID
func (tx *Transaction)setTXID() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	CheckErr("setTXID", err)

	hash := sha256.Sum256(buffer.Bytes())
	tx.TXID = hash[:]
}

func (tx *Transaction) IsCoinBase() bool {
	if len(tx.TXInputs) == 1{
		if len(tx.TXInputs[0].TXID) == 0 &&tx.TXInputs[0].Vout == -1{
			return true
		}
	}
	return false
}

// NewCoinbaseTx : initial block without input
func NewCoinbaseTx(address string, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("reword to %s %f btc", address, reward)
	}
	input := TXInput{ []byte{}, -1, data}
	output := TXOutput{reward, address}
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{output}}
	tx.setTXID()
	return &tx
}