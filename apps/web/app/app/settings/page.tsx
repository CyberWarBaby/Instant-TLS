'use client'

import { useEffect, useState } from 'react'
import { User as UserIcon, Mail, Calendar, Shield, ArrowRight } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { api, User } from '@/lib/api'
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

export default function SettingsPage() {
  const [user, setUser] = useState<User | null>(null)

  useEffect(() => {
    api.getUser().then(setUser).catch(console.error)
  }, [])

  const formatDate = (date: string) => {
    return new Date(date).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    })
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold">Settings</h1>
        <p className="text-muted-foreground mt-1">Manage your account settings</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <UserIcon className="h-5 w-5" />
            Account Information
          </CardTitle>
          <CardDescription>Your account details</CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          {user ? (
            <>
              <div className="flex items-center gap-4">
                <div className="h-16 w-16 bg-primary/10 rounded-full flex items-center justify-center">
                  <span className="text-2xl text-primary font-bold">
                    {user.email.charAt(0).toUpperCase()}
                  </span>
                </div>
                <div>
                  <p className="font-medium">{user.email}</p>
                  <PlanBadge plan={user.plan} />
                </div>
              </div>

              <div className="grid gap-4 pt-4 border-t">
                <div className="flex items-center gap-3">
                  <Mail className="h-5 w-5 text-muted-foreground" />
                  <div>
                    <p className="text-sm text-muted-foreground">Email</p>
                    <p className="font-medium">{user.email}</p>
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <Shield className="h-5 w-5 text-muted-foreground" />
                  <div>
                    <p className="text-sm text-muted-foreground">Plan</p>
                    <p className="font-medium capitalize">{user.plan}</p>
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <Calendar className="h-5 w-5 text-muted-foreground" />
                  <div>
                    <p className="text-sm text-muted-foreground">Member since</p>
                    <p className="font-medium">{formatDate(user.created_at)}</p>
                  </div>
                </div>
              </div>
            </>
          ) : (
            <div className="flex items-center justify-center py-8">
              <div className="animate-spin h-6 w-6 border-2 border-primary border-t-transparent rounded-full" />
            </div>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Plan Details</CardTitle>
              <CardDescription>Your current subscription</CardDescription>
            </div>
            <Link href="/app/pricing">
              <Button variant="outline" size="sm">
                Manage Plan
                <ArrowRight className="h-4 w-4 ml-2" />
              </Button>
            </Link>
          </div>
        </CardHeader>
        <CardContent>
          {user && (
            <div className="space-y-4">
              <Link href="/app/pricing" className="block">
                <div className="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors cursor-pointer">
                  <div>
                    <p className="font-medium capitalize">{user.plan} Plan</p>
                    <p className="text-sm text-muted-foreground">
                      {user.plan === 'free'
                        ? '1 wildcard certificate'
                        : user.plan === 'pro'
                        ? 'Unlimited wildcard certificates'
                        : 'Unlimited certificates + Team features'}
                    </p>
                  </div>
                  <div className="flex items-center gap-3">
                    <PlanBadge plan={user.plan} />
                    <ArrowRight className="h-4 w-4 text-muted-foreground" />
                  </div>
                </div>
              </Link>

              {user.plan === 'pro' && (
                <div className="p-4 bg-purple-50 border border-purple-100 rounded-lg">
                  <div className="flex items-center gap-2 mb-2">
                    <Shield className="h-5 w-5 text-purple-600" />
                    <p className="font-medium text-purple-900">Pro Plan Active</p>
                  </div>
                  <ul className="text-sm text-purple-700 space-y-1">
                    <li>✓ Unlimited wildcard certificates</li>
                    <li>✓ Priority support</li>
                    <li>✓ All CLI features</li>
                    <li>✓ HTTPS reverse proxy (instanttls serve)</li>
                  </ul>
                </div>
              )}

              {user.plan === 'team' && (
                <div className="p-4 bg-blue-50 border border-blue-100 rounded-lg">
                  <div className="flex items-center gap-2 mb-2">
                    <Shield className="h-5 w-5 text-blue-600" />
                    <p className="font-medium text-blue-900">Team Plan Active</p>
                  </div>
                  <ul className="text-sm text-blue-700 space-y-1">
                    <li>✓ Everything in Pro</li>
                    <li>✓ Team member management</li>
                    <li>✓ Shared certificates</li>
                    <li>✓ Admin dashboard</li>
                  </ul>
                </div>
              )}

              {user.plan === 'free' && (
                <Link href="/app/pricing" className="block">
                  <div className="p-4 bg-gradient-to-r from-purple-50 to-blue-50 border border-purple-200 rounded-lg hover:shadow-md transition-all cursor-pointer group">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="font-medium text-purple-900">🚀 Upgrade to Pro</p>
                        <p className="text-sm text-purple-700 mt-1">
                          Get unlimited certificates and priority support
                        </p>
                      </div>
                      <Button size="sm" className="bg-purple-600 hover:bg-purple-700 group-hover:translate-x-1 transition-transform">
                        Upgrade
                        <ArrowRight className="h-4 w-4 ml-1" />
                      </Button>
                    </div>
                  </div>
                </Link>
              )}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
