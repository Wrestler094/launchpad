-- Database initialization script for Launchpad

-- Create database (if not exists)
CREATE DATABASE IF NOT EXISTS launchpad;

-- Use the database
\c launchpad;

-- Table for storing user information
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    wallet_address VARCHAR(42) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table for storing token information
CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    total_supply BIGINT NOT NULL,
    contract_address VARCHAR(42) UNIQUE,
    owner_address VARCHAR(42) NOT NULL,
    deployment_tx_hash VARCHAR(66),
    deployed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table for storing presale information
CREATE TABLE IF NOT EXISTS presales (
    id SERIAL PRIMARY KEY,
    token_id INTEGER NOT NULL REFERENCES tokens(id),
    contract_address VARCHAR(42) UNIQUE,
    rate BIGINT NOT NULL, -- tokens per ETH
    hard_cap BIGINT NOT NULL, -- max ETH to raise (in wei)
    soft_cap BIGINT NOT NULL, -- min ETH to raise (in wei)
    duration_days INTEGER NOT NULL,
    deadline TIMESTAMP NOT NULL,
    total_raised BIGINT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    is_ended BOOLEAN DEFAULT false,
    goal_reached BOOLEAN DEFAULT false,
    deployment_tx_hash VARCHAR(66),
    landing_page_id UUID DEFAULT gen_random_uuid(),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table for tracking presale participations (optional for analytics)
CREATE TABLE IF NOT EXISTS presale_participations (
    id SERIAL PRIMARY KEY,
    presale_id INTEGER NOT NULL REFERENCES presales(id),
    participant_address VARCHAR(42) NOT NULL,
    amount_eth BIGINT NOT NULL, -- in wei
    amount_tokens BIGINT NOT NULL,
    tx_hash VARCHAR(66),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_users_wallet_address ON users(wallet_address);
CREATE INDEX IF NOT EXISTS idx_tokens_owner_address ON tokens(owner_address);
CREATE INDEX IF NOT EXISTS idx_tokens_contract_address ON tokens(contract_address);
CREATE INDEX IF NOT EXISTS idx_presales_token_id ON presales(token_id);
CREATE INDEX IF NOT EXISTS idx_presales_landing_page_id ON presales(landing_page_id);
CREATE INDEX IF NOT EXISTS idx_presale_participations_presale_id ON presale_participations(presale_id);
CREATE INDEX IF NOT EXISTS idx_presale_participations_participant ON presale_participations(participant_address);

-- Function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers to automatically update updated_at columns
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tokens_updated_at BEFORE UPDATE ON tokens FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_presales_updated_at BEFORE UPDATE ON presales FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();