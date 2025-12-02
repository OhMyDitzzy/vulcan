package core

import (
	"fmt"
	"sync"
	"github.com/OhMyDitzzy/vulcan/store"
)

type Blockchain struct {
	blocks    []*Block
	store     store.Store
	utxoSet   *UTXOSet
	mu        sync.RWMutex
	height    uint64
}

func NewBlockchain(store store.Store, utxoSet *UTXOSet) *Blockchain {
	return &Blockchain{
		blocks:  make([]*Block, 0),
		store:   store,
		utxoSet: utxoSet,
	}
}

func (bc *Blockchain) Initialize() error {
	height, err := bc.store.GetHeight()
	if err != nil || height == 0 {
		return bc.createGenesisBlock()
	}
	return bc.loadFromStore()
}

func (bc *Blockchain) createGenesisBlock() error {
	genesis := NewGenesisBlock()
	bc.blocks = append(bc.blocks, genesis)
	bc.height = 0
	
	for _, tx := range genesis.Transactions {
		bc.utxoSet.ApplyTransaction(tx)
	}
	
	data, err := genesis.ToJSON()
	if err != nil {
		return err
	}
	return bc.store.SaveBlock(genesis.Index, genesis.Hash, data)
}

func (bc *Blockchain) AddBlock(block *Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	
	if err := bc.ValidateBlock(block); err != nil {
		return fmt.Errorf("invalid block: %w", err)
	}
	
	for _, tx := range block.Transactions {
		if err := bc.utxoSet.ApplyTransaction(tx); err != nil {
			return fmt.Errorf("failed to apply transaction: %w", err)
		}
	}
	
	bc.blocks = append(bc.blocks, block)
	bc.height++
	
	data, err := block.ToJSON()
	if err != nil {
		return err
	}
	return bc.store.SaveBlock(block.Index, block.Hash, data)
}

func (bc *Blockchain) ValidateBlock(block *Block) error {
	if bc.height > 0 {
		lastBlock := bc.blocks[len(bc.blocks)-1]
		if block.PreviousHash != lastBlock.Hash {
			return fmt.Errorf("previous hash mismatch")
		}
	}
	
	if block.Index != bc.height+1 {
		return fmt.Errorf("invalid block index")
	}
	
	return block.Validate()
}

func (bc *Blockchain) GetHeight() uint64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.height
}

func (bc *Blockchain) GetLatestBlock() *Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	if len(bc.blocks) == 0 {
		return nil
	}
	return bc.blocks[len(bc.blocks)-1]
}

func (bc *Blockchain) GetBlock(index uint64) *Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	if index >= uint64(len(bc.blocks)) {
		return nil
	}
	return bc.blocks[index]
}

func (bc *Blockchain) GetBlockByHash(hash string) *Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	for _, block := range bc.blocks {
		if block.Hash == hash {
			return block
		}
	}
	return nil
}

func (bc *Blockchain) GetBlocks(start, limit uint64) []*Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	
	end := start + limit
	if end > uint64(len(bc.blocks)) {
		end = uint64(len(bc.blocks))
	}
	
	return bc.blocks[start:end]
}

func (bc *Blockchain) loadFromStore() error {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	
	height, err := bc.store.GetHeight()
	if err != nil {
		return err
	}
	
	for i := uint64(0); i <= height; i++ {
		data, err := bc.store.GetBlock(i)
		if err != nil {
			return fmt.Errorf("failed to load block %d: %w", i, err)
		}
		
		block, err := BlockFromJSON(data)
		if err != nil {
			return fmt.Errorf("failed to deserialize block %d: %w", i, err)
		}
		
		bc.blocks = append(bc.blocks, block)
		
		// Apply transactions to UTXO set
		for _, tx := range block.Transactions {
			bc.utxoSet.ApplyTransaction(tx)
		}
	}
	
	bc.height = height
	return nil
}