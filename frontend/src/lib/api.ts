import { 
  ApiResponse, 
  NonceResponse, 
  LoginResponse, 
  CreateTokenResponse,
  CreatePresaleResponse,
  Token,
  Presale,
  PresaleData,
  ParticipateResponse
} from '@/types'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api'

export class ApiClient {
  private baseUrl: string
  private token: string | null = null

  constructor() {
    this.baseUrl = API_BASE_URL
    this.token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null
  }

  setToken(token: string) {
    this.token = token
    if (typeof window !== 'undefined') {
      localStorage.setItem('auth_token', token)
    }
  }

  clearToken() {
    this.token = null
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_token')
    }
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string>),
    }

    if (this.token) {
      headers.Authorization = `Bearer ${this.token}`
    }

    const response = await fetch(url, {
      ...options,
      headers,
    })

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Unknown error' }))
      throw new Error(error.error || `HTTP error! status: ${response.status}`)
    }

    return response.json()
  }

  // Auth methods
  async generateNonce(address: string) {
    return this.request<ApiResponse<NonceResponse>>(`/auth/nonce?address=${address}`)
  }

  async login(address: string, signature: string, nonce: string) {
    return this.request<ApiResponse<LoginResponse>>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ address, signature, nonce }),
    })
  }

  async verifyToken() {
    return this.request<ApiResponse<{ address: string }>>('/auth/verify', {
      method: 'POST',
    })
  }

  // Token methods
  async createToken(name: string, symbol: string, totalSupply: string) {
    return this.request<ApiResponse<CreateTokenResponse>>('/token/create', {
      method: 'POST',
      body: JSON.stringify({ name, symbol, total_supply: totalSupply }),
    })
  }

  async getToken(address: string) {
    return this.request<ApiResponse<Token>>(`/token/${address}`)
  }

  async listTokens() {
    return this.request<ApiResponse<Token[]>>('/token/list')
  }

  // Presale methods
  async createPresale(tokenAddress: string, rate: string, softCap: string, hardCap: string, deadline: string) {
    return this.request<ApiResponse<CreatePresaleResponse>>('/presale/create', {
      method: 'POST',
      body: JSON.stringify({ token_address: tokenAddress, rate, soft_cap: softCap, hard_cap: hardCap, deadline }),
    })
  }

  async getPresale(id: number) {
    return this.request<ApiResponse<Presale>>(`/presale/${id}`)
  }

  async getPublicPresale(id: number) {
    return this.request<ApiResponse<PresaleData>>(`/public/presale/${id}`)
  }

  async listPresales() {
    return this.request<ApiResponse<Presale[]>>('/presale/list')
  }

  async participateInPresale(id: number, amountETH: string, txHash: string) {
    return this.request<ApiResponse<ParticipateResponse>>(`/presale/${id}/participate`, {
      method: 'POST',
      body: JSON.stringify({ amount_eth: amountETH, tx_hash: txHash }),
    })
  }
}

export const apiClient = new ApiClient()