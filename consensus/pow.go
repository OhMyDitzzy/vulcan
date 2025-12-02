package consensus

import (
	"fmt"
	"strings"
	"time"

	"github.com/OhMyDitzzy/vulcan/core"
)

// ProofOfWork implements the Proof-of-Work consensus algorithm.
// In our blockchain, miners must find a nonce that produces a block hash
// with a specific number of leading zeros (difficulty). This ensures
// that blocks are mined at a predictable rate and provides security
// against attacks by making chain rewriting computationally expensive.
type ProofOfWork struct {
	difficulty      int
	targetBlockTime time.Duration
}

// NewProofOfWork creates a new ProofOfWork instance.
// Configure the difficulty (number of leading zeros required)
// and target block time for dynamic difficulty adjustment.
func NewProofOfWork(difficulty int, targetBlockTime time.Duration) *ProofOfWork {
	if difficulty < 1 {
		difficulty = 1
	}
	if targetBlockTime == 0 {
		targetBlockTime = 10 * time.Second
	}
	
	return &ProofOfWork{
		difficulty:      difficulty,
		targetBlockTime: targetBlockTime,
	}
}

// Mine attempts to find a valid nonce for the block.
// We increment the nonce and compute the hash repeatedly until we find
// a hash that satisfies the difficulty requirement (has the required
// number of leading zeros). This is the core of the mining process.
func (pow *ProofOfWork) Mine(block *core.Block) error {
	fmt.Printf("Mining block %d with difficulty %d...\n", block.Index, pow.difficulty)
	
	startTime := time.Now()
	target := pow.getTarget()
	
	var hashesComputed uint64
	for {
		block.SetHash()
		
		hashesComputed++

		if pow.isValidHash(block.Hash, target) {
			duration := time.Since(startTime)
			hashRate := float64(hashesComputed) / duration.Seconds()
			fmt.Printf("Block %d mined! Hash: %s (took %v, %0.0f H/s)\n",
				block.Index, block.Hash, duration, hashRate)
			return nil
		}

		block.Nonce++

		if hashesComputed%100000 == 0 {
			fmt.Printf("Mining progress: %d hashes computed...\n", hashesComputed)
		}
	}
}

// getTarget returns the target string (required prefix of zeros).
// We build a string of zeros based on the difficulty level.
func (pow *ProofOfWork) getTarget() string {
	return strings.Repeat("0", pow.difficulty)
}

// isValidHash checks if a hash meets the difficulty requirement.
// Verify that the hash starts with the required number of leading zeros.
func (pow *ProofOfWork) isValidHash(hash, target string) bool {
	if len(hash) < len(target) {
		return false
	}
	return hash[:len(target)] == target
}

// ValidateBlock verifies that a block has valid Proof-of-Work.
// Check that the block's hash has the required number of leading zeros
// and that the hash is correctly computed from the block data.
func (pow *ProofOfWork) ValidateBlock(block *core.Block) error {
	expectedHash := block.ComputeHash()
	if block.Hash != expectedHash {
		return fmt.Errorf("block hash is incorrect: expected %s, got %s", expectedHash, block.Hash)
	}

	target := strings.Repeat("0", pow.difficulty)
	if !pow.isValidHash(block.Hash, target) {
		return fmt.Errorf("block hash does not meet difficulty requirement (need %d leading zeros)", pow.difficulty)
	}
	
	return nil
}

// AdjustDifficulty dynamically adjusts the mining difficulty based on recent block times.
// Increase difficulty if blocks are being mined too fast, and decrease it
// if blocks are taking too long. This helps maintain a consistent block time.
func (pow *ProofOfWork) AdjustDifficulty(recentBlocks []*core.Block) {
	if len(recentBlocks) < 10 {
		return // Need at least 10 blocks to adjust
	}

	var totalTime time.Duration
	for i := 1; i < len(recentBlocks); i++ {
		timeDiff := recentBlocks[i].Timestamp.Sub(recentBlocks[i-1].Timestamp)
		totalTime += timeDiff
	}
	avgTime := totalTime / time.Duration(len(recentBlocks)-1)

	if avgTime < pow.targetBlockTime/2 {
		pow.difficulty++
		fmt.Printf("Difficulty increased to %d (avg block time: %v)\n", pow.difficulty, avgTime)
	} else if avgTime > pow.targetBlockTime*2 && pow.difficulty > 1 {
		pow.difficulty--
		fmt.Printf("Difficulty decreased to %d (avg block time: %v)\n", pow.difficulty, avgTime)
	}
}


func (pow *ProofOfWork) GetDifficulty() int {
	return pow.difficulty
}

// SetDifficulty manually sets the mining difficulty.
// for testing or manual adjustment.
func (pow *ProofOfWork) SetDifficulty(difficulty int) {
	if difficulty < 1 {
		difficulty = 1
	}
	pow.difficulty = difficulty
}

// EstimateHashRate estimates the network hash rate based on a block.
// Calculate this from the difficulty and the time it took to mine the block.
func (pow *ProofOfWork) EstimateHashRate(block *core.Block, prevBlock *core.Block) float64 {
	if prevBlock == nil {
		return 0
	}
	
	timeDiff := block.Timestamp.Sub(prevBlock.Timestamp).Seconds()
	if timeDiff == 0 {
		return 0
	}
	
	// Approximate number of hashes needed: 16^difficulty
	targetHashes := 1.0
	for i := 0; i < pow.difficulty; i++ {
		targetHashes *= 16
	}
	
	hashRate := targetHashes / timeDiff
	return hashRate
}