package services

import (
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wrestler094/launchpad/internal/contracts"
	"github.com/wrestler094/launchpad/internal/storage"
)

// PresaleService handles presale-related operations
type PresaleService struct {
	client *contracts.Client
	db     *sql.DB
}

// CreatePresaleRequest represents a presale creation request
type CreatePresaleRequest struct {
	TokenAddress string `json:"token_address"`
	Rate         string `json:"rate"`
	SoftCap      string `json:"soft_cap"`
	HardCap      string `json:"hard_cap"`
	Deadline     string `json:"deadline"` // ISO 8601 format
}

// CreatePresaleResponse represents a presale creation response
type CreatePresaleResponse struct {
	Address    string           `json:"address"`
	TxHash     string           `json:"tx_hash"`
	LandingURL string           `json:"landing_url"`
	Presale    *storage.Presale `json:"presale"`
}

// ParticipateRequest represents a presale participation request
type ParticipateRequest struct {
	AmountETH string `json:"amount_eth"`
	TxHash    string `json:"tx_hash"`
}

// ParticipateResponse represents a presale participation response
type ParticipateResponse struct {
	AmountTokens   string `json:"amount_tokens"`
	Participation  *storage.PresaleParticipation `json:"participation"`
}

// NewPresaleService creates a new presale service
func NewPresaleService(client *contracts.Client, db *sql.DB) *PresaleService {
	return &PresaleService{
		client: client,
		db:     db,
	}
}

// CreatePresale creates a new presale
func (p *PresaleService) CreatePresale(creatorAddress string, req *CreatePresaleRequest) (*CreatePresaleResponse, error) {
	// Validate input
	if !common.IsHexAddress(creatorAddress) {
		return nil, fmt.Errorf("invalid creator address")
	}

	if !common.IsHexAddress(req.TokenAddress) {
		return nil, fmt.Errorf("invalid token address")
	}

	// Parse numbers
	rate, ok := new(big.Int).SetString(req.Rate, 10)
	if !ok {
		return nil, fmt.Errorf("invalid rate")
	}

	softCap, ok := new(big.Int).SetString(req.SoftCap, 10)
	if !ok {
		return nil, fmt.Errorf("invalid soft cap")
	}

	hardCap, ok := new(big.Int).SetString(req.HardCap, 10)
	if !ok {
		return nil, fmt.Errorf("invalid hard cap")
	}

	// Parse deadline
	deadline, err := time.Parse(time.RFC3339, req.Deadline)
	if err != nil {
		return nil, fmt.Errorf("invalid deadline format: %w", err)
	}

	if deadline.Before(time.Now()) {
		return nil, fmt.Errorf("deadline must be in the future")
	}

	// Verify token exists
	tokenExists, err := p.verifyTokenExists(req.TokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	if !tokenExists {
		return nil, fmt.Errorf("token not found")
	}

	// For now, simulate presale creation
	// In a real implementation, you would call the factory contract
	
	// Generate mock address and tx hash
	presaleAddress := common.HexToAddress("0x" + fmt.Sprintf("%040x", len(req.TokenAddress)+int(softCap.Int64())))
	txHash := "0x" + fmt.Sprintf("%064x", len(req.TokenAddress)*int(rate.Int64()))

	// Store presale in database
	presale := &storage.Presale{
		Address:        presaleAddress.Hex(),
		TokenAddress:   req.TokenAddress,
		CreatorAddress: creatorAddress,
		Rate:           rate.String(),
		SoftCap:        softCap.String(),
		HardCap:        hardCap.String(),
		Deadline:       deadline,
		TxHash:         txHash,
		Active:         true,
		Finalized:      false,
	}

	err = p.storePresale(presale)
	if err != nil {
		return nil, fmt.Errorf("failed to store presale: %w", err)
	}

	// Generate landing URL
	landingURL := fmt.Sprintf("http://localhost:3000/presale/%d", presale.ID)

	return &CreatePresaleResponse{
		Address:    presale.Address,
		TxHash:     presale.TxHash,
		LandingURL: landingURL,
		Presale:    presale,
	}, nil
}

// GetPresale gets a presale by ID
func (p *PresaleService) GetPresale(id int) (*storage.Presale, error) {
	query := `
		SELECT id, address, token_address, creator_address, rate, soft_cap, hard_cap, 
		       deadline, tx_hash, active, finalized, created_at
		FROM presales 
		WHERE id = $1
	`

	presale := &storage.Presale{}
	err := p.db.QueryRow(query, id).Scan(
		&presale.ID,
		&presale.Address,
		&presale.TokenAddress,
		&presale.CreatorAddress,
		&presale.Rate,
		&presale.SoftCap,
		&presale.HardCap,
		&presale.Deadline,
		&presale.TxHash,
		&presale.Active,
		&presale.Finalized,
		&presale.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("presale not found")
		}
		return nil, fmt.Errorf("failed to get presale: %w", err)
	}

	return presale, nil
}

// ListPresales lists presales created by a user
func (p *PresaleService) ListPresales(creatorAddress string) ([]*storage.Presale, error) {
	if !common.IsHexAddress(creatorAddress) {
		return nil, fmt.Errorf("invalid creator address")
	}

	query := `
		SELECT id, address, token_address, creator_address, rate, soft_cap, hard_cap, 
		       deadline, tx_hash, active, finalized, created_at
		FROM presales 
		WHERE creator_address = $1 
		ORDER BY created_at DESC
	`

	rows, err := p.db.Query(query, creatorAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to list presales: %w", err)
	}
	defer rows.Close()

	var presales []*storage.Presale
	for rows.Next() {
		presale := &storage.Presale{}
		err := rows.Scan(
			&presale.ID,
			&presale.Address,
			&presale.TokenAddress,
			&presale.CreatorAddress,
			&presale.Rate,
			&presale.SoftCap,
			&presale.HardCap,
			&presale.Deadline,
			&presale.TxHash,
			&presale.Active,
			&presale.Finalized,
			&presale.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan presale: %w", err)
		}
		presales = append(presales, presale)
	}

	return presales, nil
}

// ParticipateInPresale records a user's participation in a presale
func (p *PresaleService) ParticipateInPresale(presaleID int, participantAddress string, req *ParticipateRequest) (*ParticipateResponse, error) {
	// Validate input
	if !common.IsHexAddress(participantAddress) {
		return nil, fmt.Errorf("invalid participant address")
	}

	amountETH, ok := new(big.Int).SetString(req.AmountETH, 10)
	if !ok {
		return nil, fmt.Errorf("invalid ETH amount")
	}

	// Get presale info
	presale, err := p.GetPresale(presaleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get presale: %w", err)
	}

	if !presale.Active {
		return nil, fmt.Errorf("presale is not active")
	}

	if time.Now().After(presale.Deadline) {
		return nil, fmt.Errorf("presale has ended")
	}

	// Calculate tokens
	rate, _ := new(big.Int).SetString(presale.Rate, 10)
	amountTokens := new(big.Int).Mul(amountETH, rate)

	// Store participation
	participation := &storage.PresaleParticipation{
		PresaleID:       presaleID,
		ParticipantAddr: participantAddress,
		AmountETH:       amountETH.String(),
		AmountTokens:    amountTokens.String(),
		TxHash:          req.TxHash,
	}

	err = p.storeParticipation(participation)
	if err != nil {
		return nil, fmt.Errorf("failed to store participation: %w", err)
	}

	return &ParticipateResponse{
		AmountTokens:  amountTokens.String(),
		Participation: participation,
	}, nil
}

// verifyTokenExists checks if a token exists in the database
func (p *PresaleService) verifyTokenExists(tokenAddress string) (bool, error) {
	query := `SELECT 1 FROM tokens WHERE address = $1`
	var exists int
	err := p.db.QueryRow(query, tokenAddress).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// storePresale stores a presale in the database
func (p *PresaleService) storePresale(presale *storage.Presale) error {
	query := `
		INSERT INTO presales (address, token_address, creator_address, rate, soft_cap, hard_cap, deadline, tx_hash, active, finalized)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at
	`

	err := p.db.QueryRow(
		query,
		presale.Address,
		presale.TokenAddress,
		presale.CreatorAddress,
		presale.Rate,
		presale.SoftCap,
		presale.HardCap,
		presale.Deadline,
		presale.TxHash,
		presale.Active,
		presale.Finalized,
	).Scan(&presale.ID, &presale.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert presale: %w", err)
	}

	return nil
}

// storeParticipation stores a participation in the database
func (p *PresaleService) storeParticipation(participation *storage.PresaleParticipation) error {
	query := `
		INSERT INTO presale_participations (presale_id, participant_address, amount_eth, amount_tokens, tx_hash)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := p.db.QueryRow(
		query,
		participation.PresaleID,
		participation.ParticipantAddr,
		participation.AmountETH,
		participation.AmountTokens,
		participation.TxHash,
	).Scan(&participation.ID, &participation.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert participation: %w", err)
	}

	return nil
}