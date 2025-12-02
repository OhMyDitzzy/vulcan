package core

import ( 
    "time"
    
    "github.com/OhMyDitzzy/vulcan/types"
)

func NewGenesisBlock() *Block {
	// Pre-funded address for testing
	preFundedAddress := "04f8a1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9"
	
	// Create coinbase transaction
	coinbase := types.NewCoinbaseTransaction(preFundedAddress, 1000000)
	
	genesis := &Block{
		Index:        0,
		Timestamp:    time.Unix(1577836800, 0), // 2020-01-01
		Transactions: []*types.Transaction{coinbase},
		Nonce:        0,
		PreviousHash: "0",
		Difficulty:   1,
	}
	
	genesis.MerkleRoot = genesis.ComputeMerkleRoot()
	genesis.Hash = genesis.ComputeHash()
	
	return genesis
}