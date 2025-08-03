package storage

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id" db:"id"`
	Address   string    `json:"address" db:"address"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Token represents a deployed token
type Token struct {
	ID             int       `json:"id" db:"id"`
	Address        string    `json:"address" db:"address"`
	Name           string    `json:"name" db:"name"`
	Symbol         string    `json:"symbol" db:"symbol"`
	TotalSupply    string    `json:"total_supply" db:"total_supply"`
	CreatorAddress string    `json:"creator_address" db:"creator_address"`
	TxHash         string    `json:"tx_hash" db:"tx_hash"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// Presale represents a token presale
type Presale struct {
	ID             int       `json:"id" db:"id"`
	Address        string    `json:"address" db:"address"`
	TokenAddress   string    `json:"token_address" db:"token_address"`
	CreatorAddress string    `json:"creator_address" db:"creator_address"`
	Rate           string    `json:"rate" db:"rate"`
	SoftCap        string    `json:"soft_cap" db:"soft_cap"`
	HardCap        string    `json:"hard_cap" db:"hard_cap"`
	Deadline       time.Time `json:"deadline" db:"deadline"`
	TxHash         string    `json:"tx_hash" db:"tx_hash"`
	Active         bool      `json:"active" db:"active"`
	Finalized      bool      `json:"finalized" db:"finalized"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// PresaleParticipation represents a user's participation in a presale
type PresaleParticipation struct {
	ID               int       `json:"id" db:"id"`
	PresaleID        int       `json:"presale_id" db:"presale_id"`
	ParticipantAddr  string    `json:"participant_address" db:"participant_address"`
	AmountETH        string    `json:"amount_eth" db:"amount_eth"`
	AmountTokens     string    `json:"amount_tokens" db:"amount_tokens"`
	TxHash           string    `json:"tx_hash" db:"tx_hash"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}