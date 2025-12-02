package api

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/OhMyDitzzy/vulcan/core"
	"github.com/OhMyDitzzy/vulcan/miner"
	"github.com/OhMyDitzzy/vulcan/p2p"
	"github.com/OhMyDitzzy/vulcan/txpool"
)

// Server represents the HTTP API server.
// we provide RESTful endpoints for interacting
// with the blockchain, managing wallets, and mining blocks.
type Server struct {
	port       int
	router     *gin.Engine
	blockchain *core.Blockchain
	mempool    *txpool.Mempool
	miner      *miner.Miner
	p2pNode    *p2p.Node
	utxoSet    *core.UTXOSet
}

// NewServer creates a new API server instance.
// initialize the Gin router with middleware and register all endpoints.
func NewServer(port int, bc *core.Blockchain, mp *txpool.Mempool, m *miner.Miner, p2p *p2p.Node, utxo *core.UTXOSet) *Server {
	gin.SetMode(gin.ReleaseMode)
	
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	router.Use(cors.New(config))
	
	router.Use(RequestLogger())
	
	server := &Server{
		port:       port,
		router:     router,
		blockchain: bc,
		mempool:    mp,
		miner:      m,
		p2pNode:    p2p,
		utxoSet:    utxo,
	}
	
	server.setupRoutes()
	return server
}

// setupRoutes registers all API endpoints.
// Organize endpoints by functionality: blockchain, wallet, transactions, mining, peers.
func (s *Server) setupRoutes() {
	api := s.router.Group("/")
	
	api.GET("/health", s.handleHealth)
	
	api.GET("/blockchain/blocks", s.handleGetBlocks)
	api.GET("/blockchain/block/:hash", s.handleGetBlock)
	api.GET("/blockchain/tx/:txid", s.handleGetTransaction)
	
	api.GET("/wallet/new", s.handleNewWallet)
	api.POST("/wallet/sign", s.handleSignTransaction)
	
	api.POST("/tx", s.handleBroadcastTransaction)
	api.GET("/mempool", s.handleGetMempool)

	api.POST("/mine", s.handleMine)

	api.GET("/balance/:address", s.handleGetBalance)
	
	api.GET("/peers", s.handleGetPeers)
	api.POST("/peers", s.handleAddPeer)
	
	api.GET("/metrics", s.handleMetrics)
}

// Start starts the API server.
//Bind to the configured port and begin serving requests.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	return s.router.Run(addr)
}