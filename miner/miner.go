package miner

import (
	"log"
	"sync"
	"time"
	"github.com/OhMyDitzzy/vulcan/consensus"
	"github.com/OhMyDitzzy/vulcan/core"
	"github.com/OhMyDitzzy/vulcan/types"
	"github.com/OhMyDitzzy/vulcan/txpool"
)

type Miner struct {
	blockchain *core.Blockchain
	mempool    *txpool.Mempool
	pow        *consensus.ProofOfWork
	utxoSet    *core.UTXOSet
	mining     bool
	mu         sync.Mutex
}

func NewMiner(bc *core.Blockchain, mp *txpool.Mempool, pow *consensus.ProofOfWork, utxo *core.UTXOSet) *Miner {
	return &Miner{
		blockchain: bc,
		mempool:    mp,
		pow:        pow,
		utxoSet:    utxo,
	}
}

func (m *Miner) Start(minerAddress string) {
	m.mu.Lock()
	m.mining = true
	m.mu.Unlock()
	
	log.Println("Miner started, waiting for transactions...")
	
	for m.IsMining() {
		if m.mempool.Size() == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		
		if err := m.MineBlock(minerAddress); err != nil {
			log.Printf("Mining failed: %v", err)
		}
		
		time.Sleep(1 * time.Second)
	}
}

func (m *Miner) Stop() {
	m.mu.Lock()
	m.mining = false
	m.mu.Unlock()
}

func (m *Miner) IsMining() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.mining
}

func (m *Miner) MineBlock(minerAddress string) error {
	txs := m.mempool.GetTransactions(100)
	
	blockReward := uint64(50)
	totalFees := uint64(0)
	for _, tx := range txs {
		totalFees += tx.Fee
	}
	
	coinbase := types.NewCoinbaseTransaction(minerAddress, blockReward+totalFees)
	allTxs := append([]*types.Transaction{coinbase}, txs...)
	
	lastBlock := m.blockchain.GetLatestBlock()
	newBlock := core.NewBlock(
		m.blockchain.GetHeight()+1,
		allTxs,
		lastBlock.Hash,
		m.pow.GetDifficulty(),
	)

	if err := m.pow.Mine(newBlock); err != nil {
		return err
	}

	if err := m.blockchain.AddBlock(newBlock); err != nil {
		return err
	}
	
	if err := m.utxoSet.Update(newBlock); err != nil {
		log.Printf("Warning: Failed to update UTXO set: %v", err)
		// We shouldn't return err here
		// UTXO will be synced on next restart
	}

	for _, tx := range txs {
		m.mempool.RemoveTransaction(tx.ID)
	}
	
	log.Printf("Block %d mined successfully! Hash: %s", newBlock.Index, newBlock.Hash)
	return nil
}