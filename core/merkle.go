package core

import (
	"crypto/sha256"
	"encoding/hex"
)

// BuildMerkleRoot constructs a Merkle tree from transaction IDs and returns the root hash.
// we recursively hash pairs of hashes until we reach a single root.
// If the number of hashes is odd, we duplicate the last hash to make it even.
// This provides an efficient way to verify transaction inclusion in a block.
func BuildMerkleRoot(txIDs []string) string {
	if len(txIDs) == 0 {
		return ""
	}
	
	leaves := make([]string, len(txIDs))
	for i, txID := range txIDs {
		leaves[i] = txID
	}
	
	return buildMerkleTree(leaves)
}

// buildMerkleTree recursively builds the Merkle tree.
// Start with the leaf nodes (transaction hashes) and combine them
// pairwise until we have a single root hash.
func buildMerkleTree(hashes []string) string {
	if len(hashes) == 1 {
		return hashes[0]
	}
	
	// If odd number of hashes, duplicate the last one
	if len(hashes)%2 != 0 {
		hashes = append(hashes, hashes[len(hashes)-1])
	}
	
	var newLevel []string
	for i := 0; i < len(hashes); i += 2 {
		combined := combineHashes(hashes[i], hashes[i+1])
		newLevel = append(newLevel, combined)
	}
	
	return buildMerkleTree(newLevel)
}

// combineHashes concatenates two hashes and returns their SHA256 hash.
// Build parent nodes in the Merkle tree from child nodes.
func combineHashes(hash1, hash2 string) string {
	combined := hash1 + hash2
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

// VerifyTransactionInclusion verifies that a transaction is included in a block
// using the Merkle proof. This allows efficient verification without needing
// all transactions in the block.
func VerifyTransactionInclusion(txID string, merkleRoot string, proof []string, index int) bool {
	currentHash := txID
	
	for i, siblingHash := range proof {
		// Determine if we should concatenate on left or right based on index
		// We use bit manipulation to determine position in the tree
		if (index>>i)&1 == 0 {
			// Transaction is on the left
			currentHash = combineHashes(currentHash, siblingHash)
		} else {
			// Transaction is on the right
			currentHash = combineHashes(siblingHash, currentHash)
		}
	}
	
	return currentHash == merkleRoot
}

// GenerateMerkleProof generates a Merkle proof for a transaction at the given index.
// The proof consists of the sibling hashes needed to reconstruct the path
// from the transaction to the root. We return this proof so that a light client
// can verify transaction inclusion without downloading the entire block.
func GenerateMerkleProof(txIDs []string, index int) []string {
	if index < 0 || index >= len(txIDs) {
		return nil
	}
	
	var proof []string
	leaves := make([]string, len(txIDs))
	copy(leaves, txIDs)
	
	currentIndex := index
	currentLevel := leaves
	
	for len(currentLevel) > 1 {
		// If odd number, duplicate last element
		if len(currentLevel)%2 != 0 {
			currentLevel = append(currentLevel, currentLevel[len(currentLevel)-1])
		}
		
		var siblingIndex int
		if currentIndex%2 == 0 {
			siblingIndex = currentIndex + 1
		} else {
			siblingIndex = currentIndex - 1
		}
		
		if siblingIndex < len(currentLevel) {
			proof = append(proof, currentLevel[siblingIndex])
		}
		
		var nextLevel []string
		for i := 0; i < len(currentLevel); i += 2 {
			combined := combineHashes(currentLevel[i], currentLevel[i+1])
			nextLevel = append(nextLevel, combined)
		}
		
		currentLevel = nextLevel
		currentIndex = currentIndex / 2
	}
	
	return proof
}