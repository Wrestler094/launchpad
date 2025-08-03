package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"launchpad-backend/middleware"
	"launchpad-backend/models"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthHandler struct {
	db        *sql.DB
	jwtSecret string
}

func NewAuthHandler(db *sql.DB, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the signature
	if !verifySignature(req.Message, req.Signature, req.WalletAddress) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	// Create or get user
	user, err := h.createOrGetUser(req.WalletAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or retrieve user"})
		return
	}

	// Generate JWT token
	token, err := h.generateJWT(req.WalletAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User:  *user,
	})
}

func (h *AuthHandler) createOrGetUser(walletAddress string) (*models.User, error) {
	var user models.User
	
	// Try to get existing user
	err := h.db.QueryRow(
		"SELECT id, wallet_address, created_at, updated_at FROM users WHERE wallet_address = $1",
		walletAddress,
	).Scan(&user.ID, &user.WalletAddress, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		// Create new user
		err = h.db.QueryRow(
			"INSERT INTO users (wallet_address) VALUES ($1) RETURNING id, wallet_address, created_at, updated_at",
			walletAddress,
		).Scan(&user.ID, &user.WalletAddress, &user.CreatedAt, &user.UpdatedAt)
		
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (h *AuthHandler) generateJWT(walletAddress string) (string, error) {
	claims := &middleware.JWTClaims{
		WalletAddress: walletAddress,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

func verifySignature(message, signature, walletAddress string) bool {
	// Remove "0x" prefix if present
	if len(signature) > 2 && signature[:2] == "0x" {
		signature = signature[2:]
	}

	// Decode the signature
	sigBytes, err := hexutil.Decode("0x" + signature)
	if err != nil {
		return false
	}

	// The signature should be 65 bytes
	if len(sigBytes) != 65 {
		return false
	}

	// Ethereum uses recovery ID of 27 or 28, but we need 0 or 1
	if sigBytes[64] >= 27 {
		sigBytes[64] -= 27
	}

	// Hash the message with Ethereum prefix
	messageHash := crypto.Keccak256Hash([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)))

	// Recover the public key
	publicKeyECDSA, err := crypto.SigToPub(messageHash[:], sigBytes)
	if err != nil {
		return false
	}

	// Get the address from the public key
	recoveredAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Compare addresses (case insensitive)
	return strings.ToLower(recoveredAddress.Hex()) == strings.ToLower(walletAddress)
}