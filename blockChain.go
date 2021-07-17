package main

import (
	"github.com/boltdb/bolt"
	"os"
)

const dbFile = "blockChain.db"
const blockBucket = "bucket"
const lastHashKey = "key"

type BlockChain struct {
	db  *bolt.DB
	tail []byte
}

// InitBlockChain CreateBlockchain creates a new blockchain DB
func InitBlockChain() *BlockChain {
	//bolt db, a key-value db. key should be the hash, and value should the serialized data
	db, err := bolt.Open(dbFile, 0600, nil)
	CheckErr("InitBlockChain 1", err)
	var lastHash []byte
	//db.ViewC)
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket != nil{
			println("the bucket already exists, please change the name")
			os.Exit(1)
		}else{
			//if there is no bucket, initialization is needed
			genesis := NewGenesisBlock()
			bucket, err := tx.CreateBucket( []byte(blockBucket))
			CheckErr( "InitBlockChain 2", err)
			err = bucket.Put(genesis.Hash, genesis.Serialize())
			CheckErr( "InitBlockChain 3", err)
			err = bucket.Put([]byte(lastHashKey), genesis.Hash)
			CheckErr( "InitBlockChain 4", err)
			lastHash = genesis.Hash
		}
		return nil
	})
	CheckErr("InitBlockChain 2", err)
	return &BlockChain{db,lastHash}
}

func getBlockChainHandler() *BlockChain {
	//bolt db, a key-value db. key should be the hash, and value should the serialized data
	db, err := bolt.Open(dbFile, 0600, nil)
	CheckErr("cannot find the database", err)
	var lastHash []byte
	//db.ViewC)
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket != nil{
			//get the hash of last block
			lastHash = bucket.Get([]byte(lastHashKey))
		}else{
			println("no bucket in the database")
			os.Exit(1)
		}
		return nil
	})
	CheckErr("InitBlockChain 2", err)
	return &BlockChain{db,lastHash}
}


func (bc *BlockChain)AddBlock(data string) {
	var prevBlockHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			os.Exit(1)
		}
		prevBlockHash = bucket.Get([]byte(lastHashKey))
		return nil
	})
	CheckErr( "AddBlock whole process 2", err)
	block := NewBlock(data, prevBlockHash)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			os.Exit(1)
		}
		err := bucket.Put(block.Hash, block.Serialize())
		CheckErr( "AddBlock1", err)
		err = bucket.Put([]byte(lastHashKey), block.Hash)
		CheckErr( "AddBlock2", err)
		bc.tail = block.Hash
		return nil
	})
	CheckErr( "AddBlock whole process 2", err)
}

type BlockChainIterator struct {
	currHash []byte
	db       *bolt.DB
}

// NewIterator create a new iterator for a given blockchain
func (bc *BlockChain)NewIterator() *BlockChainIterator{
	return &BlockChainIterator{currHash: bc.tail, db:bc.db}
}

// Next iteration process
func (it *BlockChainIterator)Next() (block *Block){
	err := it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil{
			return nil
		}

		data := bucket.Get(it.currHash)
		block = Deserialize(data)
		it.currHash = block.PreviousBlockHash
		return nil
	})
	CheckErr("Next()", err)
	return

}