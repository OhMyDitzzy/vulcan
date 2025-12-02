package txpool

import (
	"fmt"
	"sort"
	"sync"
	"github.com/OhMyDitzzy/vulcan/types"
)

type Mempool struct {
	transactions map[string]*types.Transaction
	mu           sync.RWMutex
}

func NewMempool() *Mempool {
	return &Mempool{
		transactions: make(map[string]*types.Transaction),
	}
}

func (mp *Mempool) AddTransaction(tx *types.Transaction) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	
	// Check if already exists
	if _, exists := mp.transactions[tx.ID]; exists {
		return fmt.Errorf("transaction already in mempool")
	}
	
	for _, existingTx := range mp.transactions {
		if existingTx.From == tx.From && existingTx.ID != tx.ID {
			// TODO: check UTXO conflicts
		}
	}
	
	mp.transactions[tx.ID] = tx
	return nil
}

func (mp *Mempool) RemoveTransaction(txID string) {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	delete(mp.transactions, txID)
}

func (mp *Mempool) GetTransactions(limit int) []*types.Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	
	txs := make([]*types.Transaction, 0, len(mp.transactions))
	for _, tx := range mp.transactions {
		txs = append(txs, tx)
	}
	
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Fee > txs[j].Fee
	})
	
	if len(txs) > limit {
		txs = txs[:limit]
	}
	
	return txs
}

func (mp *Mempool) GetTransaction(txID string) *types.Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	return mp.transactions[txID]
}

func (mp *Mempool) Size() int {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	return len(mp.transactions)
}

func (mp *Mempool) Clear() {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.transactions = make(map[string]*types.Transaction)
}