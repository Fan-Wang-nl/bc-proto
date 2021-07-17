package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

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
	//unlock script, specify who and how to use it
	ScriptSig string
}

type TXOutput struct {
	//payment value
	Value        float64
	//lock script, specify the address of the receiver
	ScriptPubKey string
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