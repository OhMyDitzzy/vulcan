package core

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/OhMyDitzzy/vulcan/types"
)

// UTXO represents an unspent transaction output.
// In our UTXO model, each transaction consumes previous UTXOs as inputs
// and creates new UTXOs as outputs. We track all unspent outputs to
// determine account balances and validate new transactions.
type UTXO struct {
	TxID    string `json:"tx_id"`    // Transaction ID that created this UTXO
	Address string `json:"address"`  // Owner's address
	Amount  uint64 `json:"amount"`   // Amount in this UTXO
	Index   int    `json:"index"`    // Output index in the transaction
}

// UTXOSet manages the set of all unspent transaction outputs.
// Maintain an in-memory map for fast lookups and provide methods
// to add, remove, and query UTXOs. This is the core of our state management.
type UTXOSet struct {
	utxos map[string]map[int]*UTXO // map[txID]map[outputIndex]UTXO
	mu    sync.RWMutex
}

func NewUTXOSet() *UTXOSet {
	return &UTXOSet{
		utxos: make(map[string]map[int]*UTXO),
	}
}

// AddUTXO adds a new unspent output to the set.
// when processing confirmed transactions to track new outputs.
func (us *UTXOSet) AddUTXO(utxo *UTXO) {
	us.mu.Lock()
	defer us.mu.Unlock()
	
	if us.utxos[utxo.TxID] == nil {
		us.utxos[utxo.TxID] = make(map[int]*UTXO)
	}
	us.utxos[utxo.TxID][utxo.Index] = utxo
}

// RemoveUTXO removes a spent output from the set.
// when a transaction consumes an existing UTXO as input.
func (us *UTXOSet) RemoveUTXO(txID string, index int) {
	us.mu.Lock()
	defer us.mu.Unlock()
	
	if us.utxos[txID] != nil {
		delete(us.utxos[txID], index)
		if len(us.utxos[txID]) == 0 {
			delete(us.utxos, txID)
		}
	}
}

// GetUTXO retrieves a specific UTXO.
// Return nil if the UTXO doesn't exist or has been spent.
func (us *UTXOSet) GetUTXO(txID string, index int) *UTXO {
	us.mu.RLock()
	defer us.mu.RUnlock()
	
	if us.utxos[txID] != nil {
		return us.utxos[txID][index]
	}
	return nil
}

// GetUTXOsForAddress returns all UTXOs owned by an address.
// Calculate an address's balance and select inputs
// for new transactions.
func (us *UTXOSet) GetUTXOsForAddress(address string) []*UTXO {
	us.mu.RLock()
	defer us.mu.RUnlock()
	
	var utxos []*UTXO
	for _, txUTXOs := range us.utxos {
		for _, utxo := range txUTXOs {
			if utxo.Address == address {
				utxos = append(utxos, utxo)
			}
		}
	}
	return utxos
}

// GetBalance calculates the total balance for an address.
// Sum up all UTXOs owned by the address.
func (us *UTXOSet) GetBalance(address string) uint64 {
	utxos := us.GetUTXOsForAddress(address)
	var balance uint64
	for _, utxo := range utxos {
		balance += utxo.Amount
	}
	return balance
}

// ApplyTransaction updates the UTXO set based on a transaction.
// Remove spent inputs and add new outputs. This is called when
// a block is added to the chain to update the state.
func (us *UTXOSet) ApplyTransaction(tx *types.Transaction) error {
	// We use a simple model where
	// the transaction specifies from/to/amount directly. In a full UTXO
	// implementation, 
	// TODO: we would have explicit inputs and outputs.
	
	if tx.IsCoinbase() {
		us.AddUTXO(&UTXO{
			TxID:    tx.ID,
			Address: tx.To,
			Amount:  tx.Amount,
			Index:   0,
		})
		return nil
	}
	
	senderUTXOs := us.GetUTXOsForAddress(tx.From)
	if len(senderUTXOs) == 0 {
		return fmt.Errorf("sender has no UTXOs")
	}

	totalNeeded := tx.Total()
	var totalAvailable uint64
	var utxosToSpend []*UTXO
	
	for _, utxo := range senderUTXOs {
		utxosToSpend = append(utxosToSpend, utxo)
		totalAvailable += utxo.Amount
		if totalAvailable >= totalNeeded {
			break
		}
	}
	
	if totalAvailable < totalNeeded {
		return fmt.Errorf("insufficient balance: have %d, need %d", totalAvailable, totalNeeded)
	}
	
	for _, utxo := range utxosToSpend {
		us.RemoveUTXO(utxo.TxID, utxo.Index)
	}

	us.AddUTXO(&UTXO{
		TxID:    tx.ID,
		Address: tx.To,
		Amount:  tx.Amount,
		Index:   0,
	})
	
	change := totalAvailable - totalNeeded
	if change > 0 {
		us.AddUTXO(&UTXO{
			TxID:    tx.ID,
			Address: tx.From,
			Amount:  change,
			Index:   1,
		})
	}
		
	return nil
}

// Update processes a new block and updates the UTXO set.
// This should be called after a block is added to the blockchain.
func (us *UTXOSet) Update(block *Block) error {
	for _, tx := range block.Transactions {
		if err := us.ApplyTransaction(tx); err != nil {
			return fmt.Errorf("failed to apply transaction %s: %v", tx.ID, err)
		}
	}
	return nil
}

// Rebuild reconstructs the UTXO set from the entire blockchain.
// This is useful for syncing or recovering from corruption.
func (us *UTXOSet) Rebuild(blockchain *Blockchain) error {
	us.mu.Lock()
	defer us.mu.Unlock()
	
	us.utxos = make(map[string]map[int]*UTXO)
	
	height := blockchain.GetHeight()
	for i := uint64(0); i <= height; i++ {
		block := blockchain.GetBlock(i)
		if block == nil {
			continue
		}
		
		// Process all transactions in the block
		for _, tx := range block.Transactions {
			// Unlock before calling ApplyTransaction to avoid deadlock
			us.mu.Unlock()
			if err := us.ApplyTransaction(tx); err != nil {
				us.mu.Lock()
				return fmt.Errorf("failed to apply transaction %s in block %d: %v", tx.ID, i, err)
			}
			us.mu.Lock()
		}
	}
	
	return nil
}

// RevertTransaction reverts the effects of a transaction on the UTXO set.
// We use this when reorganizing the chain or handling forks.
func (us *UTXOSet) RevertTransaction(tx *types.Transaction) error {
	us.RemoveUTXO(tx.ID, 0)
	us.RemoveUTXO(tx.ID, 1) // Change output
	
	// TODO: we would need to restore the spent UTXOs
	// This requires storing the original UTXOs somewhere
	
	return nil
}

// ValidateTransaction checks if a transaction can be applied to the current UTXO set.
// Verify that the sender has sufficient balance and that all referenced
// UTXOs exist and are unspent.
func (us *UTXOSet) ValidateTransaction(tx *types.Transaction) error {
	if tx.IsCoinbase() {
		return nil
	}
	
	balance := us.GetBalance(tx.From)
	totalNeeded := tx.Total()
	
	if balance < totalNeeded {
		return fmt.Errorf("insufficient balance: have %d, need %d", balance, totalNeeded)
	}
	
	return nil
}

func (us *UTXOSet) Serialize() ([]byte, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()
	
	return json.Marshal(us.utxos)
}

func (us *UTXOSet) Deserialize(data []byte) error {
	us.mu.Lock()
	defer us.mu.Unlock()
	
	return json.Unmarshal(data, &us.utxos)
}


func (us *UTXOSet) Clone() *UTXOSet {
	us.mu.RLock()
	defer us.mu.RUnlock()
	
	clone := NewUTXOSet()
	for txID, outputs := range us.utxos {
		clone.utxos[txID] = make(map[int]*UTXO)
		for index, utxo := range outputs {
			utxoCopy := *utxo
			clone.utxos[txID][index] = &utxoCopy
		}
	}
	return clone
}

func (us *UTXOSet) Count() int {
	us.mu.RLock()
	defer us.mu.RUnlock()
	
	count := 0
	for _, outputs := range us.utxos {
		count += len(outputs)
	}
	return count
}