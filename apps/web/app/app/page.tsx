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
    <Button variant="outline" size="sm" onClick={handleCopy} className="gap-2">
      {copied ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
      {label || (copied ? 'Copied!' : 'Copy')}
    </Button>
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
                Build from source (binary releases coming soon)
              </p>
              <div className="bg-gray-900 text-gray-100 p-4 rounded-lg font-mono text-sm overflow-x-auto">
                <div className="flex items-center justify-between">
                  <code>go install github.com/CyberWarBaby/Instant-TLS/cli/cmd/instanttls@latest</code>
                  <CopyButton text="go install github.com/CyberWarBaby/Instant-TLS/cli/cmd/instanttls@latest" />
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
              <h3 className="font-medium">Login with your token</h3>
              <p className="text-sm text-muted-foreground mb-3">
                Authenticate the CLI with your Personal Access Token
              </p>
              <div className="bg-gray-900 text-gray-100 p-4 rounded-lg font-mono text-sm">
                <div className="flex items-center justify-between">
                  <code>instanttls login</code>
                  <CopyButton text="instanttls login" />
                </div>
              </div>
            </div>
          </div>

          {/* Step 4 */}
          <div className="flex gap-4">
            <div className="h-8 w-8 rounded-full bg-muted flex items-center justify-center text-sm font-medium">
              4
            </div>
            <div className="flex-1">
              <h3 className="font-medium">Initialize your local CA</h3>
              <p className="text-sm text-muted-foreground mb-3">
                Generate a trusted Certificate Authority for your machine
              </p>
              <div className="bg-gray-900 text-gray-100 p-4 rounded-lg font-mono text-sm">
                <div className="flex items-center justify-between">
                  <code>instanttls init</code>
                  <CopyButton text="instanttls init" />
                </div>
              </div>
            </div>
          </div>

          {/* Step 5 */}
          <div className="flex gap-4">
            <div className="h-8 w-8 rounded-full bg-muted flex items-center justify-center text-sm font-medium">
              5
            </div>
            <div className="flex-1">
              <h3 className="font-medium">Generate your first certificate</h3>
              <p className="text-sm text-muted-foreground mb-3">
                Create a wildcard certificate for local development
              </p>
              <div className="bg-gray-900 text-gray-100 p-4 rounded-lg font-mono text-sm">
                <div className="flex items-center justify-between">
                  <code>instanttls cert "*.local.test"</code>
                  <CopyButton text='instanttls cert "*.local.test"' />
                </div>
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
