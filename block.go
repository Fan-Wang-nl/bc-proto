package main

//import "crypto/sha256"

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
}

func NewBlock(data string, PreviousBlockHash []byte) *Block {
	var block Block

	//sha256.Sum256()
	return &block
}
