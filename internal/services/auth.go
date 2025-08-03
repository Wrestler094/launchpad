package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v5"
)

// AuthService handles authentication
type AuthService struct {
	jwtSecret []byte
	nonces    map[string]NonceInfo // In production, use Redis
}

// NonceInfo stores nonce information
type NonceInfo struct {
	Nonce     string
	ExpiresAt time.Time
}

// LoginRequest represents a login request
type LoginRequest struct {
	Address   string `json:"address"`
	Signature string `json:"signature"`
	Nonce     string `json:"nonce"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token   string `json:"token"`
	Address string `json:"address"`
}

// NonceResponse represents a nonce response
type NonceResponse struct {
	Nonce string `json:"nonce"`
}

// Claims represents JWT claims
type Claims struct {
	Address string `json:"address"`
	jwt.RegisteredClaims
}

// NewAuthService creates a new auth service
func NewAuthService() *AuthService {
	// In production, load this from environment
	jwtSecret := []byte("your-secret-key-change-this-in-production")
	
	return &AuthService{
		jwtSecret: jwtSecret,
		nonces:    make(map[string]NonceInfo),
	}
}

// GenerateNonce generates a new nonce for an address
func (a *AuthService) GenerateNonce(address string) (*NonceResponse, error) {
	// Validate address
	if !common.IsHexAddress(address) {
		return nil, fmt.Errorf("invalid ethereum address")
	}

	// Generate random nonce
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	
	nonce := hex.EncodeToString(bytes)
	
	// Store nonce with expiration (5 minutes)
	a.nonces[address] = NonceInfo{
		Nonce:     nonce,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	return &NonceResponse{
		Nonce: nonce,
	}, nil
}

// Login authenticates a user with MetaMask signature
func (a *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// Validate address
	if !common.IsHexAddress(req.Address) {
		return nil, fmt.Errorf("invalid ethereum address")
	}

	// Check nonce
	nonceInfo, exists := a.nonces[req.Address]
	if !exists {
		return nil, fmt.Errorf("nonce not found")
	}

	if time.Now().After(nonceInfo.ExpiresAt) {
		delete(a.nonces, req.Address)
		return nil, fmt.Errorf("nonce expired")
	}

	if nonceInfo.Nonce != req.Nonce {
		return nil, fmt.Errorf("invalid nonce")
	}

	// Verify signature
	message := fmt.Sprintf("Sign this message to authenticate with Launchpad.\n\nNonce: %s", req.Nonce)
	
	if !a.verifySignature(req.Address, message, req.Signature) {
		return nil, fmt.Errorf("invalid signature")
	}

	// Clean up nonce
	delete(a.nonces, req.Address)

	// Generate JWT token
	token, err := a.generateJWT(req.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		Token:   token,
		Address: req.Address,
	}, nil
}

// VerifyToken verifies a JWT token and returns the address
func (a *AuthService) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// generateJWT generates a JWT token for an address
func (a *AuthService) generateJWT(address string) (string, error) {
	claims := &Claims{
		Address: address,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtSecret)
}

// verifySignature verifies an Ethereum signature
func (a *AuthService) verifySignature(address, message, signature string) bool {
	// Hash the message
	hash := accounts.TextHash([]byte(message))

	// Decode signature
	sig, err := hex.DecodeString(signature[2:]) // Remove 0x prefix
	if err != nil {
		return false
	}

	// Adjust recovery ID for Ethereum
	if sig[64] == 27 || sig[64] == 28 {
		sig[64] -= 27
	}

	// Recover public key
	pubKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return false
	}

	// Get address from public key
	recoveredAddress := crypto.PubkeyToAddress(*pubKey)

	// Compare addresses
	return recoveredAddress.Hex() == address
}