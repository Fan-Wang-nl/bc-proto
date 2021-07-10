package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

type ProofOfWork struct{
	block *Block
	target *big.Int
}

const targetBits =24
func NewProofOfWork(block *Block) *ProofOfWork {
	//000000000000000.... 01
	target := big.NewInt( 1)
	//0x000000100000000000000
	target.Lsh(target, uint(256 - targetBits))
	pow := ProofOfWork{block:block, target: target}
	return &pow
}

func (pow *ProofOfWork)PrepareData(nonce int64)[]byte{
	block := pow.block
	tmp := [][]byte{
		IntToByte(block.Version),
		block.PreviousBlockHash,
		block.MerkelRoot,
		IntToByte(block.TimeStamp),
		IntToByte(targetBits),
		IntToByte(block.Nonce),
		block.Data}
	data := bytes.Join(tmp, []byte{})
	return data
}


func (pow *ProofOfWork)Run() (int64, []byte){
	var hash [32]byte
	var nonce int64 = 0
	var hashInt big.Int
	println("begin mining")
	fmt.Printf("target hash: %x\n", pow.target.Bytes())
	for nonce < math.MaxInt64 {
		data := pow.PrepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		// Cmp compares x and y and returns:
		//-1 if x < y
		//0 if x == y
		//+1 if x > y
		//
		if hashInt.Cmp(pow.target) == -1 {
			//fmt.Printf("found nonce, nonce :%d, hash : %x\n", nonce, hash)
			fmt.Printf("found hash : %x\n, nonce :%d, ", hash, nonce)
			break
		} else {
			//fmt.Printf("not found nonce, current nonce :%d, hash : %x\n", nonce, hash)
			nonce++

		}
	}
	return nonce, hash[:]

}