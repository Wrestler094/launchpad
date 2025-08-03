package models

import (
	"time"
	"database/sql"
)

type User struct {
	ID            int       `json:"id" db:"id"`
	WalletAddress string    `json:"wallet_address" db:"wallet_address"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type Token struct {
	ID               int            `json:"id" db:"id"`
	Name             string         `json:"name" db:"name"`
	Symbol           string         `json:"symbol" db:"symbol"`
	TotalSupply      int64          `json:"total_supply" db:"total_supply"`
	ContractAddress  sql.NullString `json:"contract_address" db:"contract_address"`
	OwnerAddress     string         `json:"owner_address" db:"owner_address"`
	DeploymentTxHash sql.NullString `json:"deployment_tx_hash" db:"deployment_tx_hash"`
	DeployedAt       sql.NullTime   `json:"deployed_at" db:"deployed_at"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
}

type Presale struct {
	ID               int            `json:"id" db:"id"`
	TokenID          int            `json:"token_id" db:"token_id"`
	ContractAddress  sql.NullString `json:"contract_address" db:"contract_address"`
	Rate             int64          `json:"rate" db:"rate"`
	HardCap          int64          `json:"hard_cap" db:"hard_cap"`
	SoftCap          int64          `json:"soft_cap" db:"soft_cap"`
	DurationDays     int            `json:"duration_days" db:"duration_days"`
	Deadline         time.Time      `json:"deadline" db:"deadline"`
	TotalRaised      int64          `json:"total_raised" db:"total_raised"`
	IsActive         bool           `json:"is_active" db:"is_active"`
	IsEnded          bool           `json:"is_ended" db:"is_ended"`
	GoalReached      bool           `json:"goal_reached" db:"goal_reached"`
	DeploymentTxHash sql.NullString `json:"deployment_tx_hash" db:"deployment_tx_hash"`
	LandingPageID    string         `json:"landing_page_id" db:"landing_page_id"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
}

type PresaleParticipation struct {
	ID                  int       `json:"id" db:"id"`
	PresaleID           int       `json:"presale_id" db:"presale_id"`
	ParticipantAddress  string    `json:"participant_address" db:"participant_address"`
	AmountETH           int64     `json:"amount_eth" db:"amount_eth"`
	AmountTokens        int64     `json:"amount_tokens" db:"amount_tokens"`
	TxHash              sql.NullString `json:"tx_hash" db:"tx_hash"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

// Request/Response DTOs
type LoginRequest struct {
	WalletAddress string `json:"wallet_address" binding:"required"`
	Signature     string `json:"signature" binding:"required"`
	Message       string `json:"message" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type TokenCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Symbol      string `json:"symbol" binding:"required"`
	TotalSupply int64  `json:"total_supply" binding:"required"`
}

type TokenCreateResponse struct {
	TokenID         int    `json:"token_id"`
	ContractAddress string `json:"contract_address"`
	TxHash          string `json:"tx_hash"`
}

type PresaleCreateRequest struct {
	TokenID      int   `json:"token_id" binding:"required"`
	Rate         int64 `json:"rate" binding:"required"`
	HardCap      int64 `json:"hard_cap" binding:"required"`
	SoftCap      int64 `json:"soft_cap" binding:"required"`
	DurationDays int   `json:"duration_days" binding:"required"`
}

type PresaleCreateResponse struct {
	PresaleID       int    `json:"presale_id"`
	ContractAddress string `json:"contract_address"`
	LandingPageURL  string `json:"landing_page_url"`
	TxHash          string `json:"tx_hash"`
}