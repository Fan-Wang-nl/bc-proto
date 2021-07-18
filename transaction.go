package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
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

func NewTransaction(){

}

//Check if the current user can consume the UTXO
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