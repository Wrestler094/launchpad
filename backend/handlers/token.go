package handlers

import (
	"database/sql"
	"net/http"

	"launchpad-backend/models"

	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	db *sql.DB
}

func NewTokenHandler(db *sql.DB) *TokenHandler {
	return &TokenHandler{db: db}
}

func (h *TokenHandler) CreateToken(c *gin.Context) {
	walletAddress, exists := c.Get("wallet_address")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.TokenCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input
	if req.TotalSupply <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Total supply must be greater than 0"})
		return
	}

	// Insert token into database
	var tokenID int
	err := h.db.QueryRow(`
		INSERT INTO tokens (name, symbol, total_supply, owner_address)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, req.Name, req.Symbol, req.TotalSupply, walletAddress).Scan(&tokenID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token record"})
		return
	}

	// TODO: Deploy contract to blockchain
	// For now, we'll just return a placeholder response
	// This should be implemented with actual contract deployment
	
	c.JSON(http.StatusOK, models.TokenCreateResponse{
		TokenID:         tokenID,
		ContractAddress: "0x" + "placeholder_address", // TODO: Real contract address
		TxHash:          "0x" + "placeholder_tx_hash",  // TODO: Real transaction hash
	})
}