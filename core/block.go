package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
	
	"github.com/OhMyDitzzy/vulcan/types"
)

// Block represents a single block in the blockchain.
// Each block contains an index, timestamp, list of transactions,
// and cryptographic links to the previous block through hashing.
// We use Proof-of-Work consensus to ensure blocks are mined securely.
type Block struct {
	Index        uint64         `json:"index"`         // Block height in the chain
	Timestamp    time.Time      `json:"timestamp"`     // Block creation time
	Transactions []*types.Transaction `json:"transactions"`  // List of transactions in this block
	Nonce        uint64         `json:"nonce"`         // Proof-of-Work nonce
	PreviousHash string         `json:"previous_hash"` // Hash of the previous block
	MerkleRoot   string         `json:"merkle_root"`   // Merkle root of all transactions
	Hash         string         `json:"hash"`          // Current block hash
	Difficulty   int            `json:"difficulty"`    // Mining difficulty (leading zeros)
}

// NewBlock creates a new block with the given parameters.
// Compute the Merkle root from the transactions to ensure
// integrity and efficient verification of transaction inclusion.
func NewBlock(index uint64, transactions []*types.Transaction, previousHash string, difficulty int) *Block {
	block := &Block{
		Index:        index,
		Timestamp:    time.Now().UTC(),
		Transactions: transactions,
		Nonce:        0,
		PreviousHash: previousHash,
		Difficulty:   difficulty,
	}
	block.MerkleRoot = block.ComputeMerkleRoot()
	return block
}

// ComputeHash calculates the SHA256 hash of the block header.
// Include all block fields in the hash to ensure tamper-proof linking.
// The hash is computed over: index, timestamp, merkle root, previous hash,
// nonce, and difficulty.
func (b *Block) ComputeHash() string {
	data := fmt.Sprintf("%d%s%s%s%d%d",
		b.Index,
		b.Timestamp.Format(time.RFC3339Nano),
		b.MerkleRoot,
		b.PreviousHash,
		b.Nonce,
		b.Difficulty,
	)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// SetHash computes and sets the block hash.
// We should call this after mining to finalize the block.
func (b *Block) SetHash() {
	b.Hash = b.ComputeHash()
}

// ComputeMerkleRoot calculates the Merkle root of all transactions in the block.
// Create a compact cryptographic commitment to all transactions,
// which allows efficient verification of transaction inclusion.
func (b *Block) ComputeMerkleRoot() string {
	if len(b.Transactions) == 0 {
		return ""
	}
	
	// Get transaction IDs
	txIDs := make([]string, len(b.Transactions))
	for i, tx := range b.Transactions {
		txIDs[i] = tx.ID
	}
	
	return BuildMerkleRoot(txIDs)
}

// Validate performs comprehensive validation on the block.
// We check block structure, hash validity, Merkle root, and all transactions.
func (b *Block) Validate() error {
	if b.Index == 0 && b.PreviousHash != "0" {
		return fmt.Errorf("genesis block must have previous hash of '0'")
	}
	
	if b.Hash == "" {
		return fmt.Errorf("block hash is empty")
	}
	
	if b.Hash != b.ComputeHash() {
		return fmt.Errorf("block hash is invalid")
	}
	
	expectedMerkleRoot := b.ComputeMerkleRoot()
	if b.MerkleRoot != expectedMerkleRoot {
		return fmt.Errorf("merkle root mismatch: expected %s, got %s", expectedMerkleRoot, b.MerkleRoot)
	}
	
	for i, tx := range b.Transactions {
		if err := tx.Validate(); err != nil {
			return fmt.Errorf("transaction %d invalid: %w", i, err)
		}
	}
	
	return nil
}

// HasValidProofOfWork checks if the block satisfies the PoW requirement.
// Verify that the block hash has the required number of leading zeros
// based on the difficulty level.
func (b *Block) HasValidProofOfWork() bool {
	requiredPrefix := ""
	for i := 0; i < b.Difficulty; i++ {
		requiredPrefix += "0"
	}
	return b.Hash[:b.Difficulty] == requiredPrefix
}

func (b *Block) ToJSON() ([]byte, error) {
	return json.Marshal(b)
}

func BlockFromJSON(data []byte) (*Block, error) {
	var block Block
	if err := json.Unmarshal(data, &block); err != nil {
		return nil, err
	}
	return &block, nil
}

// GetTransactionByID searches for a transaction in the block by its ID.
// Return the transaction if found, nil otherwise.
func (b *Block) GetTransactionByID(txID string) *types.Transaction {
	for _, tx := range b.Transactions {
		if tx.ID == txID {
			return tx
		}
	}
	return nil
}

func (b *Block) TotalFees() uint64 {
	var total uint64
	for _, tx := range b.Transactions {
		if !tx.IsCoinbase() {
			total += tx.Fee
		}
	}
	return total
}
.
func (b *Block) Size() int {
	data, err := b.ToJSON()
	if err != nil {
		return 0
	}
	return len(data)
}