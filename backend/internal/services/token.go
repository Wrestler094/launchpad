package services

import (
	"database/sql"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wrestler094/launchpad/internal/contracts"
	"github.com/wrestler094/launchpad/internal/storage"
)

// TokenService handles token-related operations
type TokenService struct {
	client *contracts.Client
	db     *sql.DB
}

// CreateTokenRequest represents a token creation request
type CreateTokenRequest struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	TotalSupply string `json:"total_supply"`
}

// CreateTokenResponse represents a token creation response
type CreateTokenResponse struct {
	Address string `json:"address"`
	TxHash  string `json:"tx_hash"`
	Token   *storage.Token `json:"token"`
}

// NewTokenService creates a new token service
func NewTokenService(client *contracts.Client, db *sql.DB) *TokenService {
	return &TokenService{
		client: client,
		db:     db,
	}
}

// CreateToken creates a new ERC20 token
func (t *TokenService) CreateToken(creatorAddress string, req *CreateTokenRequest) (*CreateTokenResponse, error) {
	// Validate input
	if req.Name == "" || req.Symbol == "" || req.TotalSupply == "" {
		return nil, fmt.Errorf("name, symbol, and total_supply are required")
	}

	if !common.IsHexAddress(creatorAddress) {
		return nil, fmt.Errorf("invalid creator address")
	}

	// Parse total supply
	totalSupply, ok := new(big.Int).SetString(req.TotalSupply, 10)
	if !ok {
		return nil, fmt.Errorf("invalid total supply")
	}

	// For now, we'll simulate token creation
	// In a real implementation, you would call the factory contract
	
	// Generate a mock address and tx hash for demonstration
	tokenAddress := common.HexToAddress("0x" + fmt.Sprintf("%040x", len(req.Name)+len(req.Symbol)))
	txHash := "0x" + fmt.Sprintf("%064x", len(req.Name)*len(req.Symbol))

	// Store token in database
	token := &storage.Token{
		Address:        tokenAddress.Hex(),
		Name:           req.Name,
		Symbol:         req.Symbol,
		TotalSupply:    totalSupply.String(),
		CreatorAddress: creatorAddress,
		TxHash:         txHash,
	}

	err := t.storeToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to store token: %w", err)
	}

	return &CreateTokenResponse{
		Address: token.Address,
		TxHash:  token.TxHash,
		Token:   token,
	}, nil
}

// GetToken gets a token by address
func (t *TokenService) GetToken(address string) (*storage.Token, error) {
	if !common.IsHexAddress(address) {
		return nil, fmt.Errorf("invalid token address")
	}

	query := `
		SELECT id, address, name, symbol, total_supply, creator_address, tx_hash, created_at
		FROM tokens 
		WHERE address = $1
	`

	token := &storage.Token{}
	err := t.db.QueryRow(query, address).Scan(
		&token.ID,
		&token.Address,
		&token.Name,
		&token.Symbol,
		&token.TotalSupply,
		&token.CreatorAddress,
		&token.TxHash,
		&token.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	return token, nil
}

// ListTokens lists tokens created by a user
func (t *TokenService) ListTokens(creatorAddress string) ([]*storage.Token, error) {
	if !common.IsHexAddress(creatorAddress) {
		return nil, fmt.Errorf("invalid creator address")
	}

	query := `
		SELECT id, address, name, symbol, total_supply, creator_address, tx_hash, created_at
		FROM tokens 
		WHERE creator_address = $1 
		ORDER BY created_at DESC
	`

	rows, err := t.db.Query(query, creatorAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to list tokens: %w", err)
	}
	defer rows.Close()

	var tokens []*storage.Token
	for rows.Next() {
		token := &storage.Token{}
		err := rows.Scan(
			&token.ID,
			&token.Address,
			&token.Name,
			&token.Symbol,
			&token.TotalSupply,
			&token.CreatorAddress,
			&token.TxHash,
			&token.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan token: %w", err)
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

// storeToken stores a token in the database
func (t *TokenService) storeToken(token *storage.Token) error {
	query := `
		INSERT INTO tokens (address, name, symbol, total_supply, creator_address, tx_hash)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := t.db.QueryRow(
		query,
		token.Address,
		token.Name,
		token.Symbol,
		token.TotalSupply,
		token.CreatorAddress,
		token.TxHash,
	).Scan(&token.ID, &token.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert token: %w", err)
	}

	return nil
}