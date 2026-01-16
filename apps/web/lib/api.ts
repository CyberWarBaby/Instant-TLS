// Use relative URL to go through Next.js API proxy
// This avoids CORS and mixed content issues when using HTTPS
const getApiUrl = () => {
  if (typeof window !== 'undefined') {
    // Client-side: use relative URL (goes through /app/api/[...path]/route.ts)
    return '/api'
  }
  // Server-side: use direct API URL
  return process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'
}

export interface User {
  id: string
  email: string
  plan: 'free' | 'pro' | 'team'
  created_at: string
}

export interface Token {
  id: string
  name: string
  prefix: string
  last_used_at: string | null
  created_at: string
}

export interface TokenCreateResponse {
  token: string
  data: Token
}

export interface AuthResponse {
  token: string
  user: User
}

class ApiClient {
  private authToken: string | null = null

  setAuthToken(token: string | null) {
    this.authToken = token
  }

  private async request<T>(method: string, path: string, body?: unknown): Promise<T> {
    const apiUrl = getApiUrl()
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    }

    if (this.authToken) {
      headers['Authorization'] = `Bearer ${this.authToken}`
    }

    const response = await fetch(`${apiUrl}${path}`, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    })

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Request failed' }))
      throw new Error(error.error || 'Request failed')
    }

    return response.json()
  }

  async register(email: string, password: string): Promise<AuthResponse> {
    return this.request('POST', '/v1/auth/register', { email, password })
  }

  async login(email: string, password: string): Promise<AuthResponse> {
    return this.request('POST', '/v1/auth/login', { email, password })
  }

  async getUser(): Promise<User> {
    return this.request('GET', '/v1/user')
  }

  async getTokens(): Promise<Token[]> {
    return this.request('GET', '/v1/tokens')
  }

  async createToken(name: string): Promise<TokenCreateResponse> {
    return this.request('POST', '/v1/tokens', { name })
  }

  async deleteToken(id: string): Promise<void> {
    return this.request('DELETE', `/v1/tokens/${id}`)
  }
}

export const api = new ApiClient()
