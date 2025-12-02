package wallet

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/OhMyDitzzy/vulcan/types"
)

// Wallet represents a user's wallet with key pair and address.
// In our blockchain, a wallet manages the user's private key and provides
// methods for signing transactions and managing funds.
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    string
}

// NewWallet creates a new wallet with a freshly generated key pair.
// Generate a random private key and derive the public key and address.
func NewWallet() (*Wallet, error) {
	privKey, err := GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}
	
	pubKey := &privKey.PublicKey
	address := PublicKeyToAddress(pubKey)
	
	return &Wallet{
		PrivateKey: privKey,
		PublicKey:  pubKey,
		Address:    address,
	}, nil
}

// FromPrivateKey creates a wallet from an existing private key hex string.
// Use this to restore wallets from backed-up private keys.
func FromPrivateKey(privKeyHex string) (*Wallet, error) {
	privKey, err := PrivateKeyFromHex(privKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}
	
	pubKey := &privKey.PublicKey
	address := PublicKeyToAddress(pubKey)
	
	return &Wallet{
		PrivateKey: privKey,
		PublicKey:  pubKey,
		Address:    address,
	}, nil
}

// SignTransaction signs a transaction with the wallet's private key.
// Compute the transaction hash and sign it, then set the signature
// on the transaction object. This proves that the wallet owner authorized
// the transaction.
func (w *Wallet) SignTransaction(tx *types.Transaction) error {
	if tx.From != w.Address {
		return fmt.Errorf("transaction sender does not match wallet address")
	}
	
	// Get data to sign
	dataToSign := tx.DataToSign()
	
	// Sign the data
	signature, err := Sign(dataToSign, w.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}
	
	// Set signature on transaction
	tx.SetSignature(signature)
	
	return nil
}

// VerifyTransactionSignature verifies that a transaction signature is valid.
// Extract the public key from the sender address and verify the signature
// against the transaction data. This ensures the transaction hasn't been
// tampered with and was actually signed by the claimed sender.
func VerifyTransactionSignature(tx *types.Transaction) (bool, error) {
	if tx.IsCoinbase() {
		return true, nil
	}
	
	pubKey, err := AddressToPublicKey(tx.From)
	if err != nil {
		return false, fmt.Errorf("invalid sender address: %w", err)
	}
	
	// Get data that was signed
	dataToSign := tx.DataToSign()

	valid, err := Verify(dataToSign, tx.Signature, pubKey)
	if err != nil {
		return false, fmt.Errorf("signature verification failed: %w", err)
	}
	
	return valid, nil
}

// Export returns the wallet's private key and address for backup..
func (w *Wallet) Export() (privateKeyHex, address string) {
	return PrivateKeyToHex(w.PrivateKey), w.Address
}

// CreateAndSignTransaction is a convenience method that creates and signs a transaction.
// Build a new transaction with the provided parameters and sign it with
// the wallet's private key in one step.
func (w *Wallet) CreateAndSignTransaction(to string, amount, fee uint64) (*types.Transaction, error) {
	tx := types.NewTransaction(w.Address, to, amount, fee)
	if err := w.SignTransaction(tx); err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}
	return tx, nil
}