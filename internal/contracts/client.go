package contracts

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client represents the blockchain client
type Client struct {
	Conn           *ethclient.Client
	PrivateKey     *ecdsa.PrivateKey
	Auth           *bind.TransactOpts
	FactoryAddress common.Address
}

// NewClient creates a new blockchain client
func NewClient() (*Client, error) {
	// Connect to local Hardhat node by default
	rpcURL := getEnv("RPC_URL", "http://localhost:8545")
	
	conn, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %w", err)
	}

	// Load private key from environment
	privateKeyHex := getEnv("PRIVATE_KEY", "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80") // Default Hardhat account
	
	privateKey, err := crypto.HexToECDSA(privateKeyHex[2:]) // Remove 0x prefix
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	// Get chain ID
	chainID, err := conn.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create auth
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth: %w", err)
	}

	// Set gas limit and gas price
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = big.NewInt(20000000000) // 20 gwei

	// Factory address (will be set after deployment)
	factoryAddressHex := getEnv("FACTORY_ADDRESS", "")
	var factoryAddress common.Address
	if factoryAddressHex != "" {
		factoryAddress = common.HexToAddress(factoryAddressHex)
	}

	return &Client{
		Conn:           conn,
		PrivateKey:     privateKey,
		Auth:           auth,
		FactoryAddress: factoryAddress,
	}, nil
}

// Close closes the blockchain connection
func (c *Client) Close() {
	c.Conn.Close()
}

// GetBalance gets the ETH balance of an address
func (c *Client) GetBalance(address common.Address) (*big.Int, error) {
	balance, err := c.Conn.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return balance, nil
}

// getEnv gets an environment variable with a fallback default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}