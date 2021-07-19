package main

import (
	"github.com/boltdb/bolt"
	"os"
)

const dbFile = "blockChain.db"
const blockBucket = "bucket"
const lastHashKey = "key"
const genesisInfo = "Initial block"

type BlockChain struct {
	db  *bolt.DB
	tail []byte
}

// InitBlockChain CreateBlockchain creates a new blockchain DB
func InitBlockChain(address string) *BlockChain {
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
			coinbase := NewCoinbaseTx(address, genesisInfo)
			var genesis = NewGenesisBlock(coinbase)
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


func (bc *BlockChain)AddBlock(txs []*Transaction) {
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
	block := NewBlock(txs, prevBlockHash)
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

//find all UTXOs of an address
func (bc *BlockChain)findUTXOTransactions(address string)  []Transaction{
	it := bc.NewIterator()
	//find all transactions which contain the output related to this address
	var UTXOTransactions []Transaction

	//find all spent transaction outputs (which are used as inputs in transactions)
	STXO := make(map[string][]int64)
	//traverse all the blocks
	for{
		block := it.Next()
		//traverse all the transactions in the block
		for _,tx := range block.Transactions{

			//traverse all inputs
			if !tx.IsCoinBase() {
				for _, input := range tx.TXInputs {
					if input.validateUTXO(address) {
						STXO[string(input.TXID)] = append(STXO[string(input.TXID)], input.Vout)
					}
				}
			}

			OUTPUTS:
			//traverse all outputs, and find all UTXO of the address
			for i, output := range tx.TXOutputs{
				//verify if this output has been consumed or not
				if STXO[string(tx.TXID)] != nil{
					indices := STXO[string(tx.TXID)]
					for _,index := range indices{
						//check with index if current output has been consumed
						if int64(i) == index{
							continue OUTPUTS
						}
					}
				}


				if output.validateUTXO(address){
					UTXOTransactions = append(UTXOTransactions, *tx)
				}
			}

		}

		if len(block.PreviousBlockHash) == 0 {
			break
		}
	}
	return UTXOTransactions
}

func (bc *BlockChain)findUTXO(address string)  []TXOutput{

	var utxos []TXOutput
	txs := bc.findUTXOTransactions(address)

	for _,tx := range txs{
		//traverse the tx
		for _,output := range tx.TXOutputs{
			if output.validateUTXO(address) {//this step seems to be redundant
				//TODO ,test pupose
				println("000 utxo value: ", output.Value, ", address: ", output.ScriptPubKey)
				utxos = append(utxos, output)
			}
		}
	}

	for _,utxo := range utxos{
		println("111 utxo value: ", utxo.Value, ", address: ", utxo.ScriptPubKey)
	}


	return utxos
}

func (bc *BlockChain) findSuitableUTXOs(address string, amount float64) (map[string][]int64, float64) {
	txs := bc.findUTXOTransactions(address)

	var total float64
	validUTXOs := make(map[string][]int64)

	FIND:
	for _,tx := range txs{
		//traverse the tx, and get the UTXO
		for index,output := range tx.TXOutputs{
			if output.validateUTXO(address) {
				if total <amount{
					total += output.Value
					validUTXOs[string(tx.TXID)] = append(validUTXOs[string(tx.TXID)], int64(index))
				}else {
					break FIND//jump out of the whole series
				}

			}
		}
	}

	return validUTXOs, total
}