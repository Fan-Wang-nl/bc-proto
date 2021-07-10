package main

import "math/big"

type ProofOfWork struct{
	block *Block
	target *big.Int
}

func NewProofOfWork() *ProofOfWork{
	
}