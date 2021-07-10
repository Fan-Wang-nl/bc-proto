package main

import (
	"bytes"
	"crypto/sha256"
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
	Data []byte
}

func NewBlock(data string, PreviousBlockHash []byte) *Block {
	var block Block

	block = Block{
		Version: 1,
		PreviousBlockHash: PreviousBlockHash,
		//Hash TODO
		MerkelRoot: []byte{},
		TimeStamp: time.Now().Unix(),
		Bits: 1,
		Nonce: 1,
		Data: []byte(data),
	}

	block.SetHash()

	return &block
}

func (block *Block)SetHash(){
	//sha256.Sum256()
	tmp := [][]byte{ 
		IntToByte(block.Version), 
		block.PreviousBlockHash, 
		block.MerkelRoot,
		IntToByte(block.TimeStamp), 
		IntToByte(block.Bits), 
		IntToByte(block.Nonce),
		block.Data} 

		data := bytes.Join(tmp, []byte{}) 
		hash := sha256.Sum256(data)
		block.Hash = hash[:]
}

func NewGenesisBlock() *Block{
	return NewBlock("Genesis Block", []byte{0})
}