'use client'

import { useEffect, useState } from 'react'
import { Copy, Check, Terminal, ArrowRight, Key } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { api, User, Token } from '@/lib/api'
import { useToast } from '@/components/ui/use-toast'
import Link from 'next/link'

function PlanBadge({ plan }: { plan: string }) {
  const colors: Record<string, string> = {
    free: 'bg-gray-100 text-gray-800',
    pro: 'bg-purple-100 text-purple-800',
    team: 'bg-blue-100 text-blue-800',
  }

  return (
    <span className={`px-3 py-1 rounded-full text-sm font-medium ${colors[plan] || colors.free}`}>
      {plan.charAt(0).toUpperCase() + plan.slice(1)}
    </span>
  )
}

function CopyButton({ text, label }: { text: string; label?: string }) {
  const [copied, setCopied] = useState(false)

  const handleCopy = () => {
    navigator.clipboard.writeText(text)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <button 
      onClick={handleCopy} 
      className="flex items-center gap-1 px-2 py-1 text-gray-300 hover:text-white hover:bg-gray-700 rounded transition-colors"
      title="Copy to clipboard"
    >
      {copied ? <Check className="h-4 w-4 text-green-400" /> : <Copy className="h-4 w-4" />}
      <span className="text-xs">{copied ? 'Copied!' : 'Copy'}</span>
    </button>
  )
}

export default function DashboardPage() {
  const [user, setUser] = useState<User | null>(null)
  const [tokens, setTokens] = useState<Token[]>([])
  const { toast } = useToast()

  useEffect(() => {
    api.getUser().then(setUser).catch(console.error)
    api.getTokens().then(setTokens).catch(console.error)
  }, [])

  const hasToken = tokens.length > 0

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground mt-1">Welcome to InstantTLS</p>
      </div>

      {/* Plan Card */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Your Plan</CardTitle>
              <CardDescription>Current subscription status</CardDescription>
            </div>
            {user && <PlanBadge plan={user.plan} />}
          </div>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div>
              {user?.plan === 'free' ? (
                <p className="text-sm text-muted-foreground">
                  Free plan includes 1 wildcard certificate.{' '}
                  <Link href="/pricing" className="text-primary hover:underline">
                    Upgrade for more
                  </Link>
                </p>
              ) : (
                <p className="text-sm text-muted-foreground">
                  You have unlimited wildcard certificates.
                </p>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Onboarding */}
      <Card>
        <CardHeader>
          <CardTitle>Get Started</CardTitle>
          <CardDescription>Follow these steps to set up InstantTLS</CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          {/* Step 1 */}
          <div className="flex gap-4">
            <div className={`h-8 w-8 rounded-full flex items-center justify-center text-sm font-medium ${hasToken ? 'bg-green-100 text-green-700' : 'bg-primary text-white'}`}>
              {hasToken ? '✓' : '1'}
            </div>
            <div className="flex-1">
              <h3 className="font-medium">Create a Personal Access Token</h3>
              <p className="text-sm text-muted-foreground mb-3">
                You'll need a token to authenticate the CLI.
              </p>
              {hasToken ? (
                <p className="text-sm text-green-600">✓ You have {tokens.length} token(s)</p>
              ) : (
                <Link href="/app/tokens">
                  <Button size="sm" className="gap-2">
                    <Key className="h-4 w-4" />
                    Create Token
                    <ArrowRight className="h-4 w-4" />
                  </Button>
                </Link>
              )}
            </div>
          </div>

          {/* Step 2 */}
          <div className="flex gap-4">
            <div className="h-8 w-8 rounded-full bg-muted flex items-center justify-center text-sm font-medium">
              2
            </div>
            <div className="flex-1">
              <h3 className="font-medium">Install the CLI</h3>
              <p className="text-sm text-muted-foreground mb-3">
                One command installs and makes it available system-wide
              </p>
              <div className="bg-gray-900 text-gray-100 p-4 rounded-lg font-mono text-sm overflow-x-auto space-y-2">
                <div className="flex items-center justify-between">
                  <code>curl -fsSL https://raw.githubusercontent.com/CyberWarBaby/Instant-TLS/main/install.sh | bash</code>
                  <CopyButton text='curl -fsSL https://raw.githubusercontent.com/CyberWarBaby/Instant-TLS/main/install.sh | bash' />
                </div>
              </div>
            </div>
          </div>

          {/* Step 3 */}
          <div className="flex gap-4">
            <div className="h-8 w-8 rounded-full bg-muted flex items-center justify-center text-sm font-medium">
              3
            </div>
            <div className="flex-1">
              <h3 className="font-medium">Run Setup</h3>
              <p className="text-sm text-muted-foreground mb-3">
                Creates CA, trusts in Chrome/Firefox, and configures everything
              </p>
              <div className="bg-gray-900 text-gray-100 p-4 rounded-lg font-mono text-sm">
                <div className="flex items-center justify-between">
                  <code>sudo instanttls setup</code>
                  <CopyButton text="sudo instanttls setup" />
                </div>
              </div>
              <p className="text-xs text-muted-foreground mt-2">
                You'll be prompted for your token during setup
              </p>
            </div>
          </div>

          {/* Step 4 */}
          <div className="flex gap-4">
            <div className="h-8 w-8 rounded-full bg-muted flex items-center justify-center text-sm font-medium">
              4
            </div>
            <div className="flex-1">
              <h3 className="font-medium">Add Domain to /etc/hosts</h3>
              <p className="text-sm text-muted-foreground mb-3">
                Point your domain to localhost
              </p>
              <div className="bg-gray-900 text-gray-100 p-4 rounded-lg font-mono text-sm">
                <div className="flex items-center justify-between">
                  <code>echo "127.0.0.1 myapp.local" | sudo tee -a /etc/hosts</code>
                  <CopyButton text='echo "127.0.0.1 myapp.local" | sudo tee -a /etc/hosts' />
                </div>
              </div>
            </div>
          </div>

          {/* Step 5 - Serve */}
          <div className="flex gap-4">
            <div className="h-8 w-8 rounded-full bg-muted flex items-center justify-center text-sm font-medium">
              5
            </div>
            <div className="flex-1">
              <h3 className="font-medium">Start HTTPS Proxy (Easiest Way!)</h3>
              <p className="text-sm text-muted-foreground mb-3">
                Run your app on HTTP, let InstantTLS handle HTTPS automatically
              </p>
              <div className="bg-gray-900 text-gray-100 p-4 rounded-lg font-mono text-sm space-y-3">
                <div>
                  <p className="text-gray-400 mb-1"># Your app runs on HTTP</p>
                  <code className="text-blue-400">node app.js</code>
                  <span className="text-gray-500 ml-2"># listening on localhost:3000</span>
                </div>
                <div>
                  <p className="text-gray-400 mb-1"># InstantTLS proxies HTTPS → HTTP</p>
                  <div className="flex items-center justify-between">
                    <code>sudo instanttls serve myapp.local --to localhost:3000</code>
                    <CopyButton text="sudo instanttls serve myapp.local --to localhost:3000" />
                  </div>
                </div>
              </div>
              <p className="text-xs text-muted-foreground mt-2">
                Now https://myapp.local works with a green lock! No code changes needed.
              </p>
            </div>
          </div>

          {/* Alternative - Manual Cert */}
          <div className="flex gap-4">
            <div className="h-8 w-8 rounded-full bg-gray-200 text-gray-500 flex items-center justify-center text-sm font-medium">
              alt
            </div>
            <div className="flex-1">
              <h3 className="font-medium text-muted-foreground">Alternative: Manual Certificate</h3>
              <p className="text-sm text-muted-foreground mb-3">
                If you prefer to configure TLS in your app directly
              </p>
              <div className="bg-gray-900 text-gray-100 p-4 rounded-lg font-mono text-sm">
                <div className="flex items-center justify-between">
                  <code>instanttls cert myapp.local</code>
                  <CopyButton text="instanttls cert myapp.local" />
                </div>
              </div>
              <div className="bg-muted p-4 rounded-lg text-sm space-y-2 mt-3">
                <p><strong>Certificate:</strong> <code>~/.instanttls/certs/myapp.local/cert.pem</code></p>
                <p><strong>Private Key:</strong> <code>~/.instanttls/certs/myapp.local/key.pem</code></p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Quick Commands */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Terminal className="h-5 w-5" />
            Quick Commands
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid md:grid-cols-2 gap-4">
            <div className="bg-gray-900 text-gray-100 p-4 rounded-lg font-mono text-sm">
              <p className="text-gray-400 mb-2"># Start HTTPS proxy</p>
              <div className="flex items-center justify-between">
                <code>instanttls serve myapp.local --to localhost:3000</code>
                <CopyButton text="instanttls serve myapp.local --to localhost:3000" />
              </div>
            </div>
            <div className="bg-gray-900 text-gray-100 p-4 rounded-lg font-mono text-sm">
              <p className="text-gray-400 mb-2"># Generate certificate</p>
              <div className="flex items-center justify-between">
                <code>instanttls cert myapp.local</code>
                <CopyButton text="instanttls cert myapp.local" />
              </div>
            </div>
            <div className="bg-gray-900 text-gray-100 p-4 rounded-lg font-mono text-sm">
              <p className="text-gray-400 mb-2"># Check your setup</p>
              <div className="flex items-center justify-between">
                <code>instanttls doctor</code>
                <CopyButton text="instanttls doctor" />
              </div>
            </div>
            <div className="bg-gray-900 text-gray-100 p-4 rounded-lg font-mono text-sm">
              <p className="text-gray-400 mb-2"># Renew expiring certs</p>
              <div className="flex items-center justify-between">
                <code>instanttls renew</code>
                <CopyButton text="instanttls renew" />
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
