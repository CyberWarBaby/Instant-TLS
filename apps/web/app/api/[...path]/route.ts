import { NextRequest, NextResponse } from 'next/server'

// Use API_URL for server-side proxy (not NEXT_PUBLIC_ since this runs on server only)
const API_BASE = process.env.NEXT_PUBLIC_API_URL || process.env.API_URL || 'http://localhost:8080'

export async function handler(req: Request) {
  const { method, headers } = req
  // Strip /api prefix from incoming path
  const url = new URL(req.url)
  const pathname = url.pathname.replace(/^\/api\//, '/')
  const targetUrl = `${API_BASE}${pathname}${url.search}`

  const init: RequestInit = {
    method,
    headers: new Headers(headers),
    credentials: 'include',
  }

  if (['POST', 'PUT', 'PATCH'].includes(method)) {
    init.body = await req.text()
  }

  const res = await fetch(targetUrl, init)
  const data = await res.text()
  return new Response(data, {
    status: res.status,
    headers: res.headers,
  })
}

// Backend expects: {"email": "", "password": ""}
// Response: {"token": "..."}
// Token is stored in localStorage as "token"
// Authorization header is set as: Authorization: Bearer <token> for all protected API calls
// Ensure login request sends correct JSON keys and stores token
// Ensure Authorization header is set for subsequent requests
// 1. Ensure login form sends correct JSON keys (likely { email, password })
// 2. On successful login, store token (e.g., in localStorage)
// 3. For authenticated requests, add Authorization: Bearer <token> header
// Example (pseudo):
// const res = await fetch('/api/v1/auth/login', { method: 'POST', body: JSON.stringify({ email, password }) })
// const { token } = await res.json();
// localStorage.setItem('token', token);
// ...
// For subsequent requests:
// fetch('/api/v1/user', { headers: { Authorization: `Bearer ${token}` } })
// ...
export { handler as GET, handler as POST, handler as PUT, handler as PATCH, handler as DELETE }
