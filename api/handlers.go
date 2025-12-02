package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/OhMyDitzzy/vulcan/types"
	"github.com/OhMyDitzzy/vulcan/wallet"
)

// handleHealth returns the health status of the node.
func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":     "healthy",
		"height":     s.blockchain.GetHeight(),
		"mempool":    s.mempool.Size(),
		"peers":      len(s.p2pNode.GetPeers()),
	})
}

// handleGetBlocks returns a paginated list of blocks.
func (s *Server) handleGetBlocks(c *gin.Context) {
	start, _ := strconv.ParseUint(c.DefaultQuery("start", "0"), 10, 64)
	limit, _ := strconv.ParseUint(c.DefaultQuery("limit", "10"), 10, 64)
	
	if limit > 100 {
		limit = 100 // Cap at 100 blocks per request
	}
	
	blocks := s.blockchain.GetBlocks(start, limit)
	
	c.JSON(http.StatusOK, gin.H{
		"blocks": blocks,
		"start":  start,
		"limit":  limit,
		"total":  s.blockchain.GetHeight() + 1,
	})
}

// handleGetBlock returns a specific block by hash.
func (s *Server) handleGetBlock(c *gin.Context) {
	hash := c.Param("hash")
	
	block := s.blockchain.GetBlockByHash(hash)
	if block == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "block not found"})
		return
	}
	
	c.JSON(http.StatusOK, block)
}

// handleGetTransaction returns a transaction by ID.
func (s *Server) handleGetTransaction(c *gin.Context) {
	txID := c.Param("txid")
	
	// Check mempool first
	tx := s.mempool.GetTransaction(txID)
	if tx != nil {
		c.JSON(http.StatusOK, gin.H{
			"transaction": tx,
			"status":      "pending",
		})
		return
	}
	
	// Search in blockchain
	height := s.blockchain.GetHeight()
	for i := uint64(0); i <= height; i++ {
		block := s.blockchain.GetBlock(i)
		if block != nil {
			tx := block.GetTransactionByID(txID)
			if tx != nil {
				c.JSON(http.StatusOK, gin.H{
					"transaction": tx,
					"status":      "confirmed",
					"block":       block.Hash,
					"block_index": block.Index,
				})
				return
			}
		}
	}
	
	c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
}

// handleNewWallet creates a new wallet and returns the keys.
func (s *Server) handleNewWallet(c *gin.Context) {
	// Require explicit consent to prevent accidental exposure of private keys
	consent := c.Query("consent")
	if consent != "true" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "You must add ?consent=true to generate a wallet",
			"warning": "This endpoint returns a private key. Never share it!",
		})
		return
	}
	
	// Generate new wallet
	w, err := wallet.NewWallet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	privateKey, address := w.Export()
	
	c.JSON(http.StatusOK, gin.H{
		"address":     address,
		"private_key": privateKey,
		"warning":     "NEVER share your private key! Store it securely offline.",
	})
}

// SignTransactionRequest represents the request to sign a transaction.
type SignTransactionRequest struct {
	PrivateKey  string               `json:"private_key" binding:"required"`
	Transaction TransactionPayload   `json:"transaction" binding:"required"`
}

// TransactionPayload represents the transaction data to sign.
type TransactionPayload struct {
	From   string `json:"from" binding:"required"`
	To     string `json:"to" binding:"required"`
	Amount uint64 `json:"amount" binding:"required"`
	Fee    uint64 `json:"fee" binding:"required"`
}

// handleSignTransaction signs a transaction with a private key.
func (s *Server) handleSignTransaction(c *gin.Context) {
	var req SignTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Create wallet from private key
	w, err := wallet.FromPrivateKey(req.PrivateKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid private key"})
		return
	}
	
	// Verify the from address matches the private key
	if w.Address != req.Transaction.From {
		c.JSON(http.StatusBadRequest, gin.H{"error": "private key does not match from address"})
		return
	}
	
	// Create and sign transaction
	tx, err := w.CreateAndSignTransaction(
		req.Transaction.To,
		req.Transaction.Amount,
		req.Transaction.Fee,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, tx)
}

// handleBroadcastTransaction broadcasts a signed transaction.
func (s *Server) handleBroadcastTransaction(c *gin.Context) {
	var tx types.Transaction
	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Validate transaction
	if err := tx.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction: " + err.Error()})
		return
	}
	
	// Verify signature
	valid, err := wallet.VerifyTransactionSignature(&tx)
	if err != nil || !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
		return
	}
	
	// Validate against UTXO set
	if err := s.utxoSet.ValidateTransaction(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "transaction validation failed: " + err.Error()})
		return
	}
	
	// Add to mempool
	if err := s.mempool.AddTransaction(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Broadcast to peers
	s.p2pNode.BroadcastTransaction(&tx)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "transaction broadcast successfully",
		"tx_id":   tx.ID,
	})
}

// handleGetMempool returns all pending transactions.
func (s *Server) handleGetMempool(c *gin.Context) {
	txs := s.mempool.GetTransactions(1000)
	
	c.JSON(http.StatusOK, gin.H{
		"transactions": txs,
		"count":        len(txs),
	})
}

// MineRequest represents a mining request.
type MineRequest struct {
	MinerAddress string `json:"miner_address" binding:"required"`
}

func (s *Server) handleMine(c *gin.Context) {
	var req MineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.MinerAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "miner_address is required"})
		return
	}

	heightBefore := s.blockchain.GetHeight()

	if err := s.miner.MineBlock(req.MinerAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "mining failed: " + err.Error()})
		return
	}
	
	latestBlock := s.blockchain.GetLatestBlock()

	if latestBlock == nil || latestBlock.Index != heightBefore + 1 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "block was not added to chain",
		})
		return
	}
	
	if s.p2pNode != nil {
		go s.p2pNode.BroadcastBlock(latestBlock)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "block mined successfully",
		"block":   latestBlock,
	})
}

// handleGetBalance returns the balance for an address.
func (s *Server) handleGetBalance(c *gin.Context) {
	address := c.Param("address")
	
	balance := s.utxoSet.GetBalance(address)
	utxos := s.utxoSet.GetUTXOsForAddress(address)
	
	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"balance": balance,
		"utxos":   utxos,
	})
}

// handleGetPeers returns the list of connected peers.
func (s *Server) handleGetPeers(c *gin.Context) {
	peers := s.p2pNode.GetPeers()
	
	c.JSON(http.StatusOK, gin.H{
		"peers": peers,
		"count": len(peers),
	})
}

// AddPeerRequest represents a request to add a peer.
type AddPeerRequest struct {
	Address string `json:"address" binding:"required"`
}

// handleAddPeer adds a new peer.
func (s *Server) handleAddPeer(c *gin.Context) {
	var req AddPeerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := s.p2pNode.AddPeer(req.Address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "peer added successfully",
		"address": req.Address,
	})
}

// handleMetrics returns Prometheus-style metrics.
func (s *Server) handleMetrics(c *gin.Context) {
	c.String(http.StatusOK, `# HELP vulcan_blockchain_height Current blockchain height
# TYPE vulcan_blockchain_height gauge
vulcan_blockchain_height %d

# HELP vulcan_mempool_size Number of transactions in mempool
# TYPE vulcan_mempool_size gauge
vulcan_mempool_size %d

# HELP vulcan_peers_count Number of connected peers
# TYPE vulcan_peers_count gauge
vulcan_peers_count %d

# HELP vulcan_utxo_count Number of unspent transaction outputs
# TYPE vulcan_utxo_count gauge
vulcan_utxo_count %d
`,
		s.blockchain.GetHeight(),
		s.mempool.Size(),
		len(s.p2pNode.GetPeers()),
		s.utxoSet.Count(),
	)
}