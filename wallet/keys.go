package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	btcecdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"
)

// GenerateKeyPair generates a new ECDSA key pair using secp256k1 curve.
// In our blockchain, we use secp256k1 (same as Bitcoin) for signatures.
// The private key must be kept secret, while the public key serves as the address.
func GenerateKeyPair() (*ecdsa.PrivateKey, error) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	
	return privKey.ToECDSA(), nil
}

// PrivateKeyToHex converts a private key to hexadecimal string.
// Export private keys for storage or transmission.
func PrivateKeyToHex(privKey *ecdsa.PrivateKey) string {
	privKeyBytes := privKey.D.Bytes()
	if len(privKeyBytes) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(privKeyBytes):], privKeyBytes)
		privKeyBytes = padded
	}
	return hex.EncodeToString(privKeyBytes)
}

// PrivateKeyFromHex reconstructs a private key from hexadecimal string.
// Import previously generated private keys.
func PrivateKeyFromHex(hexKey string) (*ecdsa.PrivateKey, error) {
	privKeyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %w", err)
	}
	
	if len(privKeyBytes) != 32 {
		return nil, fmt.Errorf("private key must be 32 bytes, got %d", len(privKeyBytes))
	}
	
	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)
	return privKey.ToECDSA(), nil
}

// PublicKeyToAddress converts a public key to an address (hex string).
// In our keys, the address is simply the uncompressed public key
// encoded as a hex string. This makes address derivation straightforward.
func PublicKeyToAddress(pubKey *ecdsa.PublicKey) string {
	// Serialize public key in uncompressed form (0x04 + X + Y)
	pubKeyBytes := append([]byte{0x04}, pubKey.X.Bytes()...)
	pubKeyBytes = append(pubKeyBytes, pubKey.Y.Bytes()...)
	return hex.EncodeToString(pubKeyBytes)
}

func AddressToPublicKey(address string) (*ecdsa.PublicKey, error) {
	pubKeyBytes, err := hex.DecodeString(address)
	if err != nil {
		return nil, fmt.Errorf("invalid address hex: %w", err)
	}
	
	if len(pubKeyBytes) != 65 {
		return nil, fmt.Errorf("invalid public key length: expected 65 bytes, got %d", len(pubKeyBytes))
	}
	
	if pubKeyBytes[0] != 0x04 {
		return nil, fmt.Errorf("invalid public key format: expected uncompressed format")
	}
	
	pubKey, err := btcec.ParsePubKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	
	return pubKey.ToECDSA(), nil
}

// Sign signs data using the private key and returns the signature as hex string.
// Deterministic ECDSA signing to ensure signature consistency.
// The signature consists of r and s values concatenated.
func Sign(data []byte, privKey *ecdsa.PrivateKey) (string, error) {
	// Convert to btcec private key for signing
	btcPrivKey, _ := btcec.PrivKeyFromBytes(privKey.D.Bytes())
	
	signature := btcecdsa.Sign(btcPrivKey, data)
	sigBytes := signature.Serialize()
	return hex.EncodeToString(sigBytes), nil
}

// Verify verifies a signature against data and public key.
// Return true if the signature is valid, false otherwise.
// This ensures that only the holder of the private key could have created the signature.
func Verify(data []byte, signature string, pubKey *ecdsa.PublicKey) (bool, error) {
	// Decode signature
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("invalid signature hex: %w", err)
	}
	
	// Parse signature 
	// try DER format first, then compact format
	sig, err := btcecdsa.ParseDERSignature(sigBytes)
	if err != nil {
		// Try parsing as compact signature (raw r + s format)
		sig, err = btcecdsa.ParseSignature(sigBytes)
		if err != nil {
			return false, fmt.Errorf("failed to parse signature: %w", err)
		}
	}

	x := pubKey.X.Bytes()
	y := pubKey.Y.Bytes()
	
	// Pad to 32 bytes if necessary
	if len(x) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(x):], x)
		x = padded
	}
	if len(y) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(y):], y)
		y = padded
	}
	
	pubKeyBytes := append([]byte{0x04}, x...)
	pubKeyBytes = append(pubKeyBytes, y...)
	
	btcPubKey, err := btcec.ParsePubKey(pubKeyBytes)
	if err != nil {
		return false, fmt.Errorf("invalid public key: %w", err)
	}
	
	// Verify signature
	valid := sig.Verify(data, btcPubKey)
	return valid, nil
}

// GenerateRandomBytes generates cryptographically secure random bytes.
// nonce generation and other cryptographic operations.
func GenerateRandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return bytes, nil
}