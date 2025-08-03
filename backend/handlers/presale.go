package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"launchpad-backend/models"

	"github.com/gin-gonic/gin"
)

type PresaleHandler struct {
	db *sql.DB
}

func NewPresaleHandler(db *sql.DB) *PresaleHandler {
	return &PresaleHandler{db: db}
}

func (h *PresaleHandler) CreatePresale(c *gin.Context) {
	walletAddress, exists := c.Get("wallet_address")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.PresaleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input
	if req.HardCap <= req.SoftCap {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hard cap must be greater than soft cap"})
		return
	}

	if req.DurationDays <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Duration must be greater than 0"})
		return
	}

	// Verify token ownership
	var tokenOwner string
	err := h.db.QueryRow("SELECT owner_address FROM tokens WHERE id = $1", req.TokenID).Scan(&tokenOwner)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify token ownership"})
		return
	}

	if tokenOwner != walletAddress {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't own this token"})
		return
	}

	// Calculate deadline
	deadline := time.Now().AddDate(0, 0, req.DurationDays)

	// Insert presale into database
	var presaleID int
	var landingPageID string
	err = h.db.QueryRow(`
		INSERT INTO presales (token_id, rate, hard_cap, soft_cap, duration_days, deadline)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, landing_page_id
	`, req.TokenID, req.Rate, req.HardCap, req.SoftCap, req.DurationDays, deadline).Scan(&presaleID, &landingPageID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create presale record"})
		return
	}

	// TODO: Deploy presale contract to blockchain
	// For now, we'll just return a placeholder response

	landingPageURL := fmt.Sprintf("http://localhost:3000/presale/%s", landingPageID)

	c.JSON(http.StatusOK, models.PresaleCreateResponse{
		PresaleID:       presaleID,
		ContractAddress: "0x" + "placeholder_presale_address", // TODO: Real contract address
		LandingPageURL:  landingPageURL,
		TxHash:          "0x" + "placeholder_presale_tx_hash", // TODO: Real transaction hash
	})
}

func (h *PresaleHandler) GetPresale(c *gin.Context) {
	landingPageID := c.Param("id")
	if landingPageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Presale ID required"})
		return
	}

	// Get presale with token information
	var presale models.Presale
	var token models.Token
	
	query := `
		SELECT 
			p.id, p.token_id, p.contract_address, p.rate, p.hard_cap, p.soft_cap,
			p.duration_days, p.deadline, p.total_raised, p.is_active, p.is_ended,
			p.goal_reached, p.deployment_tx_hash, p.landing_page_id, p.created_at, p.updated_at,
			t.name, t.symbol, t.total_supply, t.contract_address as token_contract_address
		FROM presales p
		JOIN tokens t ON p.token_id = t.id
		WHERE p.landing_page_id = $1
	`

	err := h.db.QueryRow(query, landingPageID).Scan(
		&presale.ID, &presale.TokenID, &presale.ContractAddress, &presale.Rate,
		&presale.HardCap, &presale.SoftCap, &presale.DurationDays, &presale.Deadline,
		&presale.TotalRaised, &presale.IsActive, &presale.IsEnded, &presale.GoalReached,
		&presale.DeploymentTxHash, &presale.LandingPageID, &presale.CreatedAt, &presale.UpdatedAt,
		&token.Name, &token.Symbol, &token.TotalSupply, &token.ContractAddress,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Presale not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve presale"})
		return
	}

	response := gin.H{
		"presale": presale,
		"token":   token,
	}

	c.JSON(http.StatusOK, response)
}

func (h *PresaleHandler) Participate(c *gin.Context) {
	presaleIDStr := c.Param("id")
	_, err := strconv.Atoi(presaleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid presale ID"})
		return
	}

	// TODO: Implement participation logic
	// This would involve:
	// 1. Validating the transaction on the blockchain
	// 2. Recording the participation in the database
	// 3. Updating the total raised amount

	c.JSON(http.StatusNotImplemented, gin.H{"error": "Participation endpoint not implemented yet"})
}