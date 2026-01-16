'use client'

import { useEffect, useState } from 'react'
import { Key, Plus, Trash2, Copy, Check, AlertCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { api, Token } from '@/lib/api'
import { useToast } from '@/components/ui/use-toast'

function CopyButton({ text }: { text: string }) {
  const [copied, setCopied] = useState(false)

  const handleCopy = () => {
    navigator.clipboard.writeText(text)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <Button variant="ghost" size="icon" onClick={handleCopy}>
      {copied ? <Check className="h-4 w-4 text-green-500" /> : <Copy className="h-4 w-4" />}
    </Button>
  )
}

export default function TokensPage() {
  const [tokens, setTokens] = useState<Token[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isCreating, setIsCreating] = useState(false)
  const [newTokenName, setNewTokenName] = useState('')
  const [createdToken, setCreatedToken] = useState<string | null>(null)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [tokenToDelete, setTokenToDelete] = useState<Token | null>(null)
  const { toast } = useToast()

  const loadTokens = async () => {
    try {
      const data = await api.getTokens()
      setTokens(data)
    } catch (error) {
      toast({
        title: 'Failed to load tokens',
        description: error instanceof Error ? error.message : 'Unknown error',
        variant: 'destructive',
      })
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    loadTokens()
  }, [])

  const handleCreate = async () => {
    if (!newTokenName.trim()) return

    setIsCreating(true)
    try {
      const response = await api.createToken(newTokenName.trim())
      setCreatedToken(response.token)
      setTokens([response.data, ...tokens])
      setNewTokenName('')
      toast({
        title: 'Token created',
        description: 'Make sure to copy your token now. You won\'t be able to see it again!',
      })
    } catch (error) {
      toast({
        title: 'Failed to create token',
        description: error instanceof Error ? error.message : 'Unknown error',
        variant: 'destructive',
      })
    } finally {
      setIsCreating(false)
    }
  }

  const handleDelete = async () => {
    if (!tokenToDelete) return

    try {
      await api.deleteToken(tokenToDelete.id)
      setTokens(tokens.filter(t => t.id !== tokenToDelete.id))
      setDeleteDialogOpen(false)
      setTokenToDelete(null)
      toast({
        title: 'Token revoked',
        description: 'The token has been permanently deleted.',
      })
    } catch (error) {
      toast({
        title: 'Failed to revoke token',
        description: error instanceof Error ? error.message : 'Unknown error',
        variant: 'destructive',
      })
    }
  }

  const closeDialog = () => {
    setDialogOpen(false)
    setCreatedToken(null)
    setNewTokenName('')
  }

  const formatDate = (date: string) => {
    return new Date(date).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    })
  }

  return (
    <div className="space-y-6 sm:space-y-8">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold">Access Tokens</h1>
          <p className="text-muted-foreground mt-1 text-sm sm:text-base">
            Manage tokens for CLI authentication
          </p>
        </div>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button className="gap-2 w-full sm:w-auto">
              <Plus className="h-4 w-4" />
              Create Token
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-[95vw] sm:max-w-lg">
            {createdToken ? (
              <>
                <DialogHeader>
                  <DialogTitle className="text-lg">Token Created</DialogTitle>
                  <DialogDescription className="text-sm">
                    Copy your token now. You won't see it again!
                  </DialogDescription>
                </DialogHeader>
                <div className="space-y-4">
                  <div className="flex items-center gap-2 p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
                    <AlertCircle className="h-5 w-5 text-yellow-600 flex-shrink-0" />
                    <p className="text-sm text-yellow-800">
                      Make sure to copy your personal access token now. You won't be able to see it again!
                    </p>
                  </div>
                  <div className="relative">
                    <div className="bg-gray-100 p-4 rounded-lg font-mono text-sm break-all">
                      {createdToken}
                    </div>
                    <div className="absolute top-2 right-2">
                      <CopyButton text={createdToken} />
                    </div>
                  </div>
                </div>
                <DialogFooter>
                  <Button onClick={closeDialog}>Done</Button>
                </DialogFooter>
              </>
            ) : (
              <>
                <DialogHeader>
                  <DialogTitle>Create Personal Access Token</DialogTitle>
                  <DialogDescription>
                    Give your token a descriptive name so you can identify it later.
                  </DialogDescription>
                </DialogHeader>
                <div className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="name">Token Name</Label>
                    <Input
                      id="name"
                      placeholder="e.g., MacBook Pro, Work Laptop"
                      value={newTokenName}
                      onChange={(e) => setNewTokenName(e.target.value)}
                      onKeyDown={(e) => e.key === 'Enter' && handleCreate()}
                    />
                  </div>
                </div>
                <DialogFooter>
                  <Button variant="outline" onClick={closeDialog}>
                    Cancel
                  </Button>
                  <Button onClick={handleCreate} disabled={isCreating || !newTokenName.trim()}>
                    {isCreating ? 'Creating...' : 'Create Token'}
                  </Button>
                </DialogFooter>
              </>
            )}
          </DialogContent>
        </Dialog>
      </div>

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Revoke Token</DialogTitle>
            <DialogDescription>
              Are you sure you want to revoke "{tokenToDelete?.name}"? This action cannot be undone. Any applications using this token will no longer be able to authenticate.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteDialogOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
              Revoke Token
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Card>
        <CardHeader className="pb-2 sm:pb-4">
          <CardTitle className="flex items-center gap-2 text-lg sm:text-xl">
            <Key className="h-4 w-4 sm:h-5 sm:w-5" />
            Your Tokens
          </CardTitle>
          <CardDescription className="text-xs sm:text-sm">
            Tokens authenticate the InstantTLS CLI
          </CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="flex items-center justify-center py-8">
              <div className="animate-spin h-6 w-6 border-2 border-primary border-t-transparent rounded-full" />
            </div>
          ) : tokens.length === 0 ? (
            <div className="text-center py-8">
              <Key className="h-10 w-10 sm:h-12 sm:w-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="font-medium mb-1 text-sm sm:text-base">No tokens yet</h3>
              <p className="text-xs sm:text-sm text-muted-foreground mb-4">
                Create a token to authenticate the CLI
              </p>
              <Button onClick={() => setDialogOpen(true)} className="gap-2" size="sm">
                <Plus className="h-4 w-4" />
                Create Token
              </Button>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b">
                    <th className="text-left py-3 px-4 font-medium">Name</th>
                    <th className="text-left py-3 px-4 font-medium">Token Prefix</th>
                    <th className="text-left py-3 px-4 font-medium">Last Used</th>
                    <th className="text-left py-3 px-4 font-medium">Created</th>
                    <th className="text-right py-3 px-4 font-medium">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {tokens.map((token) => (
                    <tr key={token.id} className="border-b last:border-0">
                      <td className="py-3 px-4 font-medium">{token.name}</td>
                      <td className="py-3 px-4">
                        <code className="text-sm bg-gray-100 px-2 py-1 rounded">
                          {token.prefix}...
                        </code>
                      </td>
                      <td className="py-3 px-4 text-muted-foreground">
                        {token.last_used_at ? formatDate(token.last_used_at) : 'Never'}
                      </td>
                      <td className="py-3 px-4 text-muted-foreground">
                        {formatDate(token.created_at)}
                      </td>
                      <td className="py-3 px-4 text-right">
                        <Button
                          variant="ghost"
                          size="icon"
                          className="text-red-500 hover:text-red-700 hover:bg-red-50"
                          onClick={() => {
                            setTokenToDelete(token)
                            setDeleteDialogOpen(true)
                          }}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
