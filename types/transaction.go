package types

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// Transaction represents a blockchain transaction with ECDSA signature.
// In our implementation, we use a UTXO model where transactions consume
// inputs and create outputs. Each transaction must be properly signed
// by the sender's private key to be considered valid.
type Transaction struct {
	ID        string    `json:"id"`         // SHA256 hash of transaction data
	From      string    `json:"from"`       // Sender's public key (hex)
	To        string    `json:"to"`         // Recipient's public key (hex)
	Amount    uint64    `json:"amount"`     // Amount to transfer
	Fee       uint64    `json:"fee"`        // Mining fee
	Signature string    `json:"signature"`  // ECDSA signature (hex)
	Timestamp time.Time `json:"timestamp"`  // Transaction creation time
}

// NewTransaction creates a new unsigned transaction.
// We must call Sign() on this transaction before broadcasting it
// to ensure authenticity and prevent tampering.
func NewTransaction(from, to string, amount, fee uint64) *Transaction {
	return &Transaction{
		From:      from,
		To:        to,
		Amount:    amount,
		Fee:       fee,
		Timestamp: time.Now().UTC(),
	}
}

// Hash computes the SHA256 hash of the transaction.
// Calculate the hash over all transaction fields except the ID itself
// to create a unique identifier for this transaction.
func (tx *Transaction) Hash() string {
	data := fmt.Sprintf("%s%s%d%d%s%s",
		tx.From,
		tx.To,
		tx.Amount,
		tx.Fee,
		tx.Signature,
		tx.Timestamp.Format(time.RFC3339Nano),
	)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// DataToSign returns the data that should be signed by the sender.
// Include all transaction fields except the signature itself
// to prevent signature malleability attacks.
func (tx *Transaction) DataToSign() []byte {
	data := fmt.Sprintf("%s%s%d%d%s",
		tx.From,
		tx.To,
		tx.Amount,
		tx.Fee,
		tx.Timestamp.Format(time.RFC3339Nano),
	)
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

// SetSignature sets the signature and computes the transaction ID.
// must call this after signing to finalize the transaction.
func (tx *Transaction) SetSignature(signature string) {
	tx.Signature = signature
	tx.ID = tx.Hash()
}

// Validate performs basic validation on the transaction.
// check that all required fields are present and have valid values.
func (tx *Transaction) Validate() error {
	if tx.IsCoinbase() {
		if tx.To == "" {
			return fmt.Errorf("to address is required")
		}
		if tx.Amount == 0 {
			return fmt.Errorf("amount must be greater than zero")
		}
		if tx.ID == "" {
			return fmt.Errorf("transaction ID must be set")
		}
		if tx.ID != tx.Hash() {
			return fmt.Errorf("transaction ID mismatch")
		}
		return nil
	}
	
	if tx.From == "" {
		return fmt.Errorf("from address is required")
	}
	if tx.To == "" {
		return fmt.Errorf("to address is required")
	}
	if tx.Amount == 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if tx.Fee == 0 {
		return fmt.Errorf("fee must be greater than zero")
	}
	if tx.Signature == "" {
		return fmt.Errorf("transaction must be signed")
	}
	if tx.ID == "" {
		return fmt.Errorf("transaction ID must be set")
	}
	if tx.ID != tx.Hash() {
		return fmt.Errorf("transaction ID mismatch")
	}
	return nil
}

// IsCoinbase returns true if this is a coinbase transaction.
// In our blockchain, coinbase transactions have empty "from" field
// and are used to reward miners for creating new blocks.
func (tx *Transaction) IsCoinbase() bool {
	return tx.From == "" && tx.Signature == "coinbase"
}

// Total returns the total amount including fee.
// calculate the total deduction from the sender's balance.
func (tx *Transaction) Total() uint64 {
	return tx.Amount + tx.Fee
}

func (tx *Transaction) ToJSON() ([]byte, error) {
	return json.Marshal(tx)
}

func FromJSON(data []byte) (*Transaction, error) {
	var tx Transaction
	if err := json.Unmarshal(data, &tx); err != nil {
		return nil, err
	}
	return &tx, nil
}

// NewCoinbaseTransaction creates a new coinbase transaction for mining rewards.
// reward the miner who successfully mines a block.
// The coinbase transaction doesn't have a sender and uses a special signature.
func NewCoinbaseTransaction(to string, amount uint64) *Transaction {
	tx := &Transaction{
		From:      "",
		To:        to,
		Amount:    amount,
		Fee:       45,
		Signature: "coinbase",
		Timestamp: time.Now().UTC(),
	}
	tx.ID = tx.Hash()
	return tx
}