import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Shield, Terminal, Zap, CheckCircle } from 'lucide-react'

export default function HomePage() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white">
      {/* Navigation */}
      <nav className="border-b bg-white/50 backdrop-blur-sm sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16 items-center">
            <div className="flex items-center gap-2">
              <Shield className="h-8 w-8 text-primary" />
              <span className="text-xl font-bold">InstantTLS</span>
            </div>
            <div className="flex items-center gap-4">
              <Link href="/login">
                <Button variant="ghost">Login</Button>
              </Link>
              <Link href="/register">
                <Button>Get Started</Button>
              </Link>
            </div>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="py-20 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-5xl font-bold tracking-tight text-gray-900 sm:text-6xl">
            Trusted HTTPS Locally
            <span className="text-primary"> with Zero Browser Warnings</span>
          </h1>
          <p className="mt-6 text-xl text-gray-600 max-w-2xl mx-auto">
            InstantTLS generates trusted local certificates for development. 
            No more clicking through security warnings.
          </p>
          <div className="mt-10 flex justify-center gap-4">
            <Link href="/register">
              <Button size="lg" className="h-12 px-8">
                Start Free
              </Button>
            </Link>
            <Link href="#install">
              <Button size="lg" variant="outline" className="h-12 px-8">
                View Install
              </Button>
            </Link>
          </div>
        </div>
      </section>

      {/* Install Section */}
      <section id="install" className="py-16 px-4 sm:px-6 lg:px-8 bg-gray-900 text-white">
        <div className="max-w-3xl mx-auto">
          <h2 className="text-2xl font-bold text-center mb-8">Install in Seconds</h2>
          <div className="bg-gray-800 rounded-lg p-6 font-mono">
            <div className="flex items-center gap-2 text-gray-400 mb-4">
              <Terminal className="h-4 w-4" />
              <span>Terminal</span>
            </div>
            <div className="space-y-2 text-green-400">
              <p><span className="text-gray-500"># Install CLI (coming soon - for now, build from source)</span></p>
              <p>$ go install github.com/CyberWarBaby/Instant-TLS/cli@latest</p>
              <p>&nbsp;</p>
              <p><span className="text-gray-500"># Login with your token</span></p>
              <p>$ instanttls login</p>
              <p>&nbsp;</p>
              <p><span className="text-gray-500"># Initialize local CA</span></p>
              <p>$ instanttls init</p>
              <p>&nbsp;</p>
              <p><span className="text-gray-500"># Generate wildcard certificate</span></p>
              <p>$ instanttls cert "*.local.test"</p>
            </div>
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="py-20 px-4 sm:px-6 lg:px-8">
        <div className="max-w-5xl mx-auto">
          <h2 className="text-3xl font-bold text-center mb-12">Everything You Need</h2>
          <div className="grid md:grid-cols-3 gap-8">
            <div className="bg-white p-6 rounded-xl border shadow-sm">
              <div className="h-12 w-12 bg-primary/10 rounded-lg flex items-center justify-center mb-4">
                <Shield className="h-6 w-6 text-primary" />
              </div>
              <h3 className="text-lg font-semibold mb-2">Local CA</h3>
              <p className="text-gray-600">
                Generate a local Certificate Authority trusted by your system and browsers.
              </p>
            </div>
            <div className="bg-white p-6 rounded-xl border shadow-sm">
              <div className="h-12 w-12 bg-primary/10 rounded-lg flex items-center justify-center mb-4">
                <Zap className="h-6 w-6 text-primary" />
              </div>
              <h3 className="text-lg font-semibold mb-2">Wildcard Certs</h3>
              <p className="text-gray-600">
                Create *.local.test certificates for all your local development domains.
              </p>
            </div>
            <div className="bg-white p-6 rounded-xl border shadow-sm">
              <div className="h-12 w-12 bg-primary/10 rounded-lg flex items-center justify-center mb-4">
                <Terminal className="h-6 w-6 text-primary" />
              </div>
              <h3 className="text-lg font-semibold mb-2">Beautiful CLI</h3>
              <p className="text-gray-600">
                Polished command-line experience with colors, spinners, and clear output.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Pricing */}
      <section className="py-20 px-4 sm:px-6 lg:px-8 bg-gray-50">
        <div className="max-w-5xl mx-auto">
          <h2 className="text-3xl font-bold text-center mb-4">Simple Pricing</h2>
          <p className="text-center text-gray-600 mb-12">Start free, upgrade when you need more.</p>
          
          <div className="grid md:grid-cols-3 gap-8">
            {/* Free */}
            <div className="bg-white p-8 rounded-xl border shadow-sm">
              <h3 className="text-lg font-semibold mb-2">Free</h3>
              <div className="text-3xl font-bold mb-4">$0<span className="text-lg font-normal text-gray-500">/mo</span></div>
              <ul className="space-y-3 mb-8">
                <li className="flex items-center gap-2">
                  <CheckCircle className="h-5 w-5 text-green-500" />
                  <span>1 wildcard certificate</span>
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle className="h-5 w-5 text-green-500" />
                  <span>Local CA generation</span>
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle className="h-5 w-5 text-green-500" />
                  <span>Auto-renew</span>
                </li>
              </ul>
              <Link href="/register">
                <Button variant="outline" className="w-full">Get Started</Button>
              </Link>
            </div>

            {/* Pro */}
            <div className="bg-white p-8 rounded-xl border-2 border-primary shadow-lg relative">
              <div className="absolute -top-3 left-1/2 -translate-x-1/2 bg-primary text-white px-3 py-1 rounded-full text-sm font-medium">
                Popular
              </div>
              <h3 className="text-lg font-semibold mb-2">Pro</h3>
              <div className="text-3xl font-bold mb-4">$9<span className="text-lg font-normal text-gray-500">/mo</span></div>
              <ul className="space-y-3 mb-8">
                <li className="flex items-center gap-2">
                  <CheckCircle className="h-5 w-5 text-green-500" />
                  <span>Unlimited certificates</span>
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle className="h-5 w-5 text-green-500" />
                  <span>Priority support</span>
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle className="h-5 w-5 text-green-500" />
                  <span>Multiple machines</span>
                </li>
              </ul>
              <Link href="/register">
                <Button className="w-full">Upgrade to Pro</Button>
              </Link>
            </div>

            {/* Team */}
            <div className="bg-white p-8 rounded-xl border shadow-sm">
              <h3 className="text-lg font-semibold mb-2">Team</h3>
              <div className="text-3xl font-bold mb-4">$29<span className="text-lg font-normal text-gray-500">/mo</span></div>
              <ul className="space-y-3 mb-8">
                <li className="flex items-center gap-2">
                  <CheckCircle className="h-5 w-5 text-green-500" />
                  <span>Everything in Pro</span>
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle className="h-5 w-5 text-green-500" />
                  <span>Team management</span>
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle className="h-5 w-5 text-green-500" />
                  <span>Shared certificates</span>
                </li>
              </ul>
              <Link href="/register">
                <Button variant="outline" className="w-full">Contact Sales</Button>
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="py-12 px-4 sm:px-6 lg:px-8 border-t">
        <div className="max-w-5xl mx-auto flex flex-col md:flex-row justify-between items-center">
          <div className="flex items-center gap-2 mb-4 md:mb-0">
            <Shield className="h-6 w-6 text-primary" />
            <span className="font-bold">InstantTLS</span>
          </div>
          <p className="text-gray-500 text-sm">
            © 2026 InstantTLS. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  )
}
