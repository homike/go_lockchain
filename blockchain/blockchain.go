package blockchain

import (
	"github.com/boltdb/bolt"
)

const dbFile = "blockchain_%s.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Blockchain struct {
	Tip []byte // 只储存顶端的块
	DB  *bolt.DB
}

// func NewBlockchain() *Blockchain {
// 	return &Blockchain{[]*Block{NewGenesisBlock()}}
// }

func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket)) // 获取储存块的桶
		if b == nil {
			genesis := NewGenesisBlock()
			b, _ := tx.CreateBucket([]byte(blocksBucket))
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("1"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("1"))
		}

		return nil
	})

	_ = err

	bc := &Blockchain{tip, db}
	return bc
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// func (bc *Blockchain) AddBlock(data string) {
// 	prevBlock := bc.Blocks[len(bc.Blocks)-1]
// 	newBlock := NewBlock(data, prevBlock.Hash)
// 	bc.Blocks = append(bc.Blocks, newBlock)
// }

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	// 获取最后一个块哈希, 使用它来挖掘一个新的块哈希
	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("1"))
		return nil
	})
	if err != nil {
		return
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		err = b.Put([]byte("1"), newBlock.Hash)
		if err != nil {
			return err
		}
		bc.Tip = newBlock.Hash

		return nil
	})
}
