'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { Check, Zap, Shield, Users, ArrowLeft, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { api, User } from '@/lib/api'
import { useToast } from '@/components/ui/use-toast'
import Link from 'next/link'

const plans = [
  {
    id: 'free',
    name: 'Free',
    price: 0,
    description: 'Perfect for trying out InstantTLS',
    features: [
      '1 wildcard certificate',
      'Local CA generation',
      'Auto-renew',
      'Community support',
    ],
    icon: Shield,
    color: 'gray',
  },
  {
    id: 'pro',
    name: 'Pro',
    price: 9,
    description: 'For professional developers',
    features: [
      'Unlimited certificates',
      'Priority support',
      'Multiple machines',
      'HTTPS reverse proxy',
      'Wildcard domains',
    ],
    icon: Zap,
    color: 'purple',
    popular: true,
  },
  {
    id: 'team',
    name: 'Team',
    price: 29,
    description: 'For development teams',
    features: [
      'Everything in Pro',
      'Team management',
      'Shared certificates',
      'Admin dashboard',
      'SSO integration',
    ],
    icon: Users,
    color: 'blue',
  },
]

export default function PricingPage() {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState<string | null>(null)
  const [showPayment, setShowPayment] = useState<string | null>(null)
  const { toast } = useToast()
  const router = useRouter()

  useEffect(() => {
    api.getUser().then(setUser).catch(console.error)
  }, [])

  const handleUpgrade = async (planId: string) => {
    if (planId === 'free') {
      // Downgrade to free
      setLoading(planId)
      try {
        await api.updatePlan(planId)
        toast({
          title: 'Plan Updated',
          description: 'You are now on the Free plan.',
        })
        setUser(prev => prev ? { ...prev, plan: 'free' } : null)
      } catch (error) {
        toast({
          title: 'Error',
          description: 'Failed to update plan. Please try again.',
          variant: 'destructive',
        })
      } finally {
        setLoading(null)
      }
    } else {
      // Show payment modal for paid plans
      setShowPayment(planId)
    }
  }

  const handlePayment = async (planId: string) => {
    setLoading(planId)
    try {
      // Simulate payment processing
      await new Promise(resolve => setTimeout(resolve, 1500))
      
      await api.updatePlan(planId)
      
      toast({
        title: '🎉 Welcome to ' + planId.charAt(0).toUpperCase() + planId.slice(1) + '!',
        description: 'Your plan has been upgraded successfully.',
      })
      
      setUser(prev => prev ? { ...prev, plan: planId as 'free' | 'pro' | 'team' } : null)
      setShowPayment(null)
    } catch (error) {
      toast({
        title: 'Payment Failed',
        description: 'Please try again or contact support.',
        variant: 'destructive',
      })
    } finally {
      setLoading(null)
    }
  }

  return (
    <div className="space-y-8">
      <div className="flex items-center gap-4">
        <Link href="/app/settings">
          <Button variant="ghost" size="icon">
            <ArrowLeft className="h-5 w-5" />
          </Button>
        </Link>
        <div>
          <h1 className="text-3xl font-bold">Pricing</h1>
          <p className="text-muted-foreground mt-1">Choose the plan that's right for you</p>
        </div>
      </div>

      <div className="grid md:grid-cols-3 gap-6">
        {plans.map((plan) => {
          const Icon = plan.icon
          const isCurrentPlan = user?.plan === plan.id
          const isPlanHigher = plans.findIndex(p => p.id === plan.id) > plans.findIndex(p => p.id === user?.plan)
          
          return (
            <Card 
              key={plan.id} 
              className={`relative ${plan.popular ? 'border-2 border-purple-500 shadow-lg' : ''}`}
            >
              {plan.popular && (
                <div className="absolute -top-3 left-1/2 -translate-x-1/2 bg-purple-500 text-white px-3 py-1 rounded-full text-sm font-medium">
                  Most Popular
                </div>
              )}
              <CardHeader>
                <div className={`h-12 w-12 rounded-lg flex items-center justify-center mb-4 ${
                  plan.color === 'gray' ? 'bg-gray-100' :
                  plan.color === 'purple' ? 'bg-purple-100' : 'bg-blue-100'
                }`}>
                  <Icon className={`h-6 w-6 ${
                    plan.color === 'gray' ? 'text-gray-600' :
                    plan.color === 'purple' ? 'text-purple-600' : 'text-blue-600'
                  }`} />
                </div>
                <CardTitle>{plan.name}</CardTitle>
                <CardDescription>{plan.description}</CardDescription>
                <div className="mt-4">
                  <span className="text-4xl font-bold">${plan.price}</span>
                  <span className="text-muted-foreground">/month</span>
                </div>
              </CardHeader>
              <CardContent className="space-y-4">
                <ul className="space-y-3">
                  {plan.features.map((feature, i) => (
                    <li key={i} className="flex items-center gap-2">
                      <Check className={`h-5 w-5 ${
                        plan.color === 'gray' ? 'text-gray-500' :
                        plan.color === 'purple' ? 'text-purple-500' : 'text-blue-500'
                      }`} />
                      <span className="text-sm">{feature}</span>
                    </li>
                  ))}
                </ul>
                
                {isCurrentPlan ? (
                  <Button variant="outline" className="w-full" disabled>
                    Current Plan
                  </Button>
                ) : (
                  <Button 
                    className={`w-full ${plan.popular ? 'bg-purple-600 hover:bg-purple-700' : ''}`}
                    variant={plan.id === 'free' ? 'outline' : 'default'}
                    onClick={() => handleUpgrade(plan.id)}
                    disabled={loading !== null}
                  >
                    {loading === plan.id ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : isPlanHigher ? (
                      `Upgrade to ${plan.name}`
                    ) : (
                      `Switch to ${plan.name}`
                    )}
                  </Button>
                )}
              </CardContent>
            </Card>
          )
        })}
      </div>

      {/* Payment Modal */}
      {showPayment && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <Card className="w-full max-w-md mx-4">
            <CardHeader>
              <CardTitle>Complete Your Upgrade</CardTitle>
              <CardDescription>
                Upgrade to {showPayment.charAt(0).toUpperCase() + showPayment.slice(1)} Plan
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="bg-gray-50 p-4 rounded-lg">
                <div className="flex justify-between items-center">
                  <span>{showPayment.charAt(0).toUpperCase() + showPayment.slice(1)} Plan</span>
                  <span className="font-bold">
                    ${plans.find(p => p.id === showPayment)?.price}/mo
                  </span>
                </div>
              </div>

              {/* Simulated Card Form */}
              <div className="space-y-4">
                <div>
                  <label className="text-sm font-medium">Card Number</label>
                  <input 
                    type="text" 
                    placeholder="4242 4242 4242 4242"
                    className="w-full mt-1 px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                    defaultValue="4242 4242 4242 4242"
                  />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="text-sm font-medium">Expiry</label>
                    <input 
                      type="text" 
                      placeholder="MM/YY"
                      className="w-full mt-1 px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                      defaultValue="12/28"
                    />
                  </div>
                  <div>
                    <label className="text-sm font-medium">CVC</label>
                    <input 
                      type="text" 
                      placeholder="123"
                      className="w-full mt-1 px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                      defaultValue="123"
                    />
                  </div>
                </div>
              </div>

              <p className="text-xs text-muted-foreground text-center">
                🔒 This is a demo. No real payment will be processed.
              </p>

              <div className="flex gap-3">
                <Button 
                  variant="outline" 
                  className="flex-1"
                  onClick={() => setShowPayment(null)}
                  disabled={loading !== null}
                >
                  Cancel
                </Button>
                <Button 
                  className="flex-1 bg-purple-600 hover:bg-purple-700"
                  onClick={() => handlePayment(showPayment)}
                  disabled={loading !== null}
                >
                  {loading ? (
                    <>
                      <Loader2 className="h-4 w-4 animate-spin mr-2" />
                      Processing...
                    </>
                  ) : (
                    `Pay $${plans.find(p => p.id === showPayment)?.price}/mo`
                  )}
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  )
}
