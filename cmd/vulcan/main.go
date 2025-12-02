package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/OhMyDitzzy/vulcan/api"
	"github.com/OhMyDitzzy/vulcan/consensus"
	"github.com/OhMyDitzzy/vulcan/core"
	"github.com/OhMyDitzzy/vulcan/miner"
	"github.com/OhMyDitzzy/vulcan/p2p"
	"github.com/OhMyDitzzy/vulcan/store"
	"github.com/OhMyDitzzy/vulcan/txpool"
)

func main() {
	// Parse command-line flags
	apiPort := flag.Int("api-port", getEnvInt("API_PORT", 8080), "API server port")
	p2pPort := flag.Int("port", getEnvInt("P2P_PORT", 6000), "P2P network port")
	dbPath := flag.String("db-path", getEnv("DB_PATH", "./data"), "Database directory path")
	peersStr := flag.String("peers", getEnv("BOOTSTRAP_PEERS", ""), "Comma-separated list of bootstrap peers")
	enableMining := flag.Bool("mining", getEnvBool("ENABLE_MINING", false), "Enable automatic mining")
	minerAddress := flag.String("miner-address", getEnv("MINER_ADDRESS", ""), "Address to receive mining rewards")
	difficulty := flag.Int("difficulty", getEnvInt("DIFFICULTY", 4), "Mining difficulty (leading zeros)")
	
	flag.Parse()

	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║    Vulcan Blockchain Node v1.0.0     ║")
	fmt.Println("╚══════════════════════════════════════╝")
	fmt.Println()

	// Initialize components
	log.Println("Initializing blockchain node...")
	
	// Create database
	db, err := store.NewBadgerStore(*dbPath)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()
	log.Printf("✓ Database initialized at %s", *dbPath)

	// Initialize UTXO set
	utxoSet := core.NewUTXOSet()
	
	// Initialize blockchain with genesis block
	blockchain := core.NewBlockchain(db, utxoSet)
	if err := blockchain.Initialize(); err != nil {
		log.Fatalf("Failed to initialize blockchain: %v", err)
	}
	log.Printf("✓ Blockchain initialized (height: %d)", blockchain.GetHeight())

	// Rebuild UTXO set from blockchain
	log.Println("Rebuilding UTXO set from blockchain...")
	if err := utxoSet.Rebuild(blockchain); err != nil {
		log.Fatalf("Failed to rebuild UTXO set: %v", err)
	}
	log.Printf("✓ UTXO set rebuilt (%d UTXOs)", utxoSet.Count())

	// Initialize transaction pool
	mempool := txpool.NewMempool()
	log.Println("✓ Transaction pool initialized")

	// Initialize consensus
	pow := consensus.NewProofOfWork(*difficulty, 10*time.Second)
	log.Printf("✓ Proof-of-Work consensus initialized (difficulty: %d)", *difficulty)

	// Initialize miner
	blockMiner := miner.NewMiner(blockchain, mempool, pow, utxoSet)
	if *enableMining {
		if *minerAddress == "" {
			log.Println("⚠ Mining enabled but no miner address specified")
		} else {
			log.Printf("✓ Miner initialized (reward address: %s)", *minerAddress)
			go blockMiner.Start(*minerAddress)
		}
	}

	// Initialize P2P network
	peers := []string{}
	if *peersStr != "" {
		peers = strings.Split(*peersStr, ",")
	}
	
	p2pNode := p2p.NewNode(*p2pPort, blockchain, mempool, peers)
	if err := p2pNode.Start(); err != nil {
		log.Fatalf("Failed to start P2P node: %v", err)
	}
	log.Printf("✓ P2P node started on port %d", *p2pPort)

	// Initialize API server
	apiServer := api.NewServer(*apiPort, blockchain, mempool, blockMiner, p2pNode, utxoSet)
	go func() {
		log.Printf("✓ API server starting on port %d", *apiPort)
		if err := apiServer.Start(); err != nil {
			log.Fatalf("Failed to start API server: %v", err)
		}
	}()

	// Print node information
	fmt.Println()
	fmt.Println("Node Information:")
	fmt.Printf("  - API Endpoint:  http://localhost:%d\n", *apiPort)
	fmt.Printf("  - P2P Address:   localhost:%d\n", *p2pPort)
	fmt.Printf("  - Blockchain Height: %d\n", blockchain.GetHeight())
	fmt.Printf("  - Total UTXOs: %d\n", utxoSet.Count())
	fmt.Printf("  - Mining: %v\n", *enableMining)
	fmt.Printf("  - Connected Peers: %d\n", len(peers))
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop the node")
	fmt.Println()

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Shutdown gracefully
	log.Println("\nShutting down node...")
	if *enableMining {
		blockMiner.Stop()
	}
	p2pNode.Stop()
	log.Println("✓ Node stopped successfully")
}

// Helper functions to read environment variables with defaults
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}