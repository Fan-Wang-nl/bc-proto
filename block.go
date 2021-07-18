package main

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Block struct {
	Version int64
	PreviousBlockHash []byte

	//to simplify the model, we create a variable of the hash of current block
	Hash []byte
	MerkelRoot []byte
	TimeStamp int64
	//indicates the level of difficulties
	Bits int64
	//random number
	Nonce int64

	//transaction information
	//Data []byte
	Transactions []*Transaction
}

func (block *Block)Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(block)
	CheckErr("Serialize", err)
	return buffer.Bytes()
}

func Deserialize(data []byte) *Block {
	if len(data) == 0 {
		return nil
	}
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	CheckErr("Deserialize", err)
	return &block
}


func NewBlock(txs []*Transaction, PreviousBlockHash []byte) *Block {
	var block Block

	block = Block{
		Version: 1,
		PreviousBlockHash: PreviousBlockHash,
		//Hash TODO
		MerkelRoot: []byte{},
		TimeStamp: time.Now().Unix(),
		Bits: targetBits,
		Nonce: 0,
		Transactions: txs}

	//block.SetHash()
	pow := NewProofOfWork(&block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash
	return &block
}



func NewGenesisBlock(coinbase *Transaction) *Block{
	return NewBlock([]*Transaction{coinbase}, []byte{})
}