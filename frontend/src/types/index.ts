export interface Token {
  id: number
  address: string
  name: string
  symbol: string
  total_supply: string
  creator_address: string
  tx_hash: string
  created_at: string
}

export interface Presale {
  id: number
  address: string
  token_address: string
  creator_address: string
  rate: string
  soft_cap: string
  hard_cap: string
  deadline: string
  tx_hash: string
  active: boolean
  finalized: boolean
  created_at: string
}

export interface PresaleParticipation {
  id: number
  presale_id: number
  participant_address: string
  amount_eth: string
  amount_tokens: string
  tx_hash: string
  created_at: string
}

export interface ApiResponse<T> {
  message: string
  data: T
}

export interface ApiError {
  error: string
}

export interface NonceResponse {
  nonce: string
}

export interface LoginResponse {
  token: string
  address: string
}

export interface CreateTokenResponse {
  address: string
  tx_hash: string
  token: Token
}

export interface CreatePresaleResponse {
  address: string
  landing_url: string
  presale: Presale
}

export interface ParticipateResponse {
  amount_tokens: string
  participation: PresaleParticipation
}

export interface PresaleData {
  presale: Presale
  token: Token
}