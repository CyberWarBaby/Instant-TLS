'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Shield, Terminal, Zap, CheckCircle, ArrowRight, Code, Server, Settings, Menu, X } from 'lucide-react'

export default function HomePage() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)

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
            
            {/* Desktop nav */}
            <div className="hidden md:flex items-center gap-4">
              <Link href="#install" className="text-sm text-gray-600 hover:text-gray-900">
                Install
              </Link>
              <Link href="#advanced" className="text-sm text-gray-600 hover:text-gray-900">
                Advanced
              </Link>
              <Link href="/login">
                <Button variant="ghost">Login</Button>
              </Link>
              <Link href="/register">
                <Button>Get Started</Button>
              </Link>
            </div>
            
            {/* Mobile menu button */}
            <button 
              className="md:hidden p-2 rounded-lg hover:bg-gray-100"
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            >
              {mobileMenuOpen ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
            </button>
          </div>
          
          {/* Mobile nav */}
          {mobileMenuOpen && (
            <div className="md:hidden py-4 border-t space-y-2">
              <Link 
                href="#install" 
                className="block px-4 py-2 text-gray-600 hover:text-gray-900 hover:bg-gray-50 rounded-lg"
                onClick={() => setMobileMenuOpen(false)}
              >
                Install
              </Link>
              <Link 
                href="#advanced" 
                className="block px-4 py-2 text-gray-600 hover:text-gray-900 hover:bg-gray-50 rounded-lg"
                onClick={() => setMobileMenuOpen(false)}
              >
                Advanced
              </Link>
              <div className="pt-2 space-y-2 px-4">
                <Link href="/login" className="block">
                  <Button variant="outline" className="w-full">Login</Button>
                </Link>
                <Link href="/register" className="block">
                  <Button className="w-full">Get Started</Button>
                </Link>
              </div>
            </div>
          )}
        </div>
      </nav>

      {/* Hero Section */}
      <section className="py-12 sm:py-20 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <div className="inline-flex items-center gap-2 bg-green-50 text-green-700 px-3 sm:px-4 py-2 rounded-full text-xs sm:text-sm font-medium mb-6">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-green-500"></span>
            </span>
            Zero code changes required
          </div>
          <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold tracking-tight text-gray-900">
            HTTPS in Seconds
            <span className="text-primary block sm:inline"> Not Hours</span>
          </h1>
          <p className="mt-6 text-lg sm:text-xl text-gray-600 max-w-2xl mx-auto px-4 sm:px-0">
            Run your app on HTTP. InstantTLS handles HTTPS automatically.
            Green lock in every browser, no configuration needed.
          </p>
          <div className="mt-10 flex flex-col sm:flex-row justify-center gap-4 px-4 sm:px-0">
            <Link href="/register">
              <Button size="lg" className="w-full sm:w-auto h-12 px-8">
                Start Free
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </Link>
            <Link href="#install">
              <Button size="lg" variant="outline" className="w-full sm:w-auto h-12 px-8">
                See How It Works
              </Button>
            </Link>
          </div>
        </div>
      </section>

      {/* Main Install Section - The Easy Way */}
      <section id="install" className="py-12 sm:py-16 px-4 sm:px-6 lg:px-8 bg-gray-900 text-white">
        <div className="max-w-4xl mx-auto">
          <div className="text-center mb-8 sm:mb-12">
            <h2 className="text-2xl sm:text-3xl font-bold mb-4">The Simplest Way to Get HTTPS Locally</h2>
            <p className="text-gray-400 text-sm sm:text-base">No nginx, no code changes, no hassle. Just works.</p>
          </div>

          <div className="grid sm:grid-cols-3 gap-4 sm:gap-6 mb-8 sm:mb-12">
            <div className="bg-gray-800/50 rounded-xl p-4 sm:p-6 border border-gray-700">
              <div className="h-10 w-10 bg-blue-500/20 rounded-lg flex items-center justify-center mb-4">
                <span className="text-blue-400 font-bold">1</span>
              </div>
              <h3 className="font-semibold mb-2">Install</h3>
              <p className="text-gray-400 text-sm">One curl command installs everything</p>
            </div>
            <div className="bg-gray-800/50 rounded-xl p-4 sm:p-6 border border-gray-700">
              <div className="h-10 w-10 bg-purple-500/20 rounded-lg flex items-center justify-center mb-4">
                <span className="text-purple-400 font-bold">2</span>
              </div>
              <h3 className="font-semibold mb-2">Setup</h3>
              <p className="text-gray-400 text-sm">Trusts CA in all your browsers automatically</p>
            </div>
            <div className="bg-gray-800/50 rounded-xl p-4 sm:p-6 border border-gray-700">
              <div className="h-10 w-10 bg-green-500/20 rounded-lg flex items-center justify-center mb-4">
                <span className="text-green-400 font-bold">3</span>
              </div>
              <h3 className="font-semibold mb-2">Serve</h3>
              <p className="text-gray-400 text-sm">HTTPS proxy to your HTTP app</p>
            </div>
          </div>

          <div className="bg-gray-800 rounded-xl p-4 sm:p-6 font-mono border border-gray-700">
            <div className="flex items-center gap-2 text-gray-400 mb-4">
              <Terminal className="h-4 w-4" />
              <span className="text-sm">Terminal</span>
            </div>
            <div className="space-y-3 text-xs sm:text-sm">
              <div>
                <p className="text-gray-500"># 1. Install (one-time)</p>
                <p className="text-green-400 break-all">$ curl -fsSL https://git.io/instanttls | bash</p>
              </div>
              <div className="pt-2">
                <p className="text-gray-500"># 2. Add domain to hosts</p>
                <p className="text-green-400 break-all">$ echo "127.0.0.1 myapp.local" | sudo tee -a /etc/hosts</p>
              </div>
              <div className="pt-2">
                <p className="text-gray-500"># 3. Run your app on HTTP</p>
                <p className="text-blue-400">$ node app.js</p>
              </div>
              <div className="pt-2">
                <p className="text-gray-500"># 4. Start HTTPS proxy</p>
                <p className="text-green-400 break-all">$ sudo instanttls serve myapp.local --to localhost:3000</p>
              </div>
              <div className="pt-4 border-t border-gray-700 mt-4">
                <p className="text-gray-400">✨ Visit <span className="text-white font-semibold">https://myapp.local</span> 🔒</p>
              </div>
            </div>
          </div>

          <div className="mt-6 sm:mt-8 text-center">
            <p className="text-gray-400 text-xs sm:text-sm">
              No sudo? Use <code className="bg-gray-800 px-2 py-1 rounded">--port 8443</code> instead
            </p>
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="py-12 sm:py-20 px-4 sm:px-6 lg:px-8">
        <div className="max-w-5xl mx-auto">
          <h2 className="text-2xl sm:text-3xl font-bold text-center mb-4">Why InstantTLS?</h2>
          <p className="text-center text-gray-600 mb-8 sm:mb-12 text-sm sm:text-base">Everything you need for local HTTPS development</p>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 sm:gap-6">
            <div className="bg-white p-4 sm:p-6 rounded-xl border shadow-sm text-center">
              <div className="h-10 w-10 sm:h-12 sm:w-12 bg-green-100 rounded-lg flex items-center justify-center mb-3 sm:mb-4 mx-auto">
                <Zap className="h-5 w-5 sm:h-6 sm:w-6 text-green-600" />
              </div>
              <h3 className="font-semibold mb-1 sm:mb-2 text-sm sm:text-base">Zero Config</h3>
              <p className="text-gray-600 text-xs sm:text-sm">
                No code changes. Your app stays on HTTP.
              </p>
            </div>
            <div className="bg-white p-4 sm:p-6 rounded-xl border shadow-sm text-center">
              <div className="h-10 w-10 sm:h-12 sm:w-12 bg-blue-100 rounded-lg flex items-center justify-center mb-3 sm:mb-4 mx-auto">
                <Shield className="h-5 w-5 sm:h-6 sm:w-6 text-blue-600" />
              </div>
              <h3 className="font-semibold mb-1 sm:mb-2 text-sm sm:text-base">Trusted CA</h3>
              <p className="text-gray-600 text-xs sm:text-sm">
                Auto-installs in Chrome, Firefox, and system.
              </p>
            </div>
            <div className="bg-white p-4 sm:p-6 rounded-xl border shadow-sm text-center">
              <div className="h-10 w-10 sm:h-12 sm:w-12 bg-purple-100 rounded-lg flex items-center justify-center mb-3 sm:mb-4 mx-auto">
                <Terminal className="h-5 w-5 sm:h-6 sm:w-6 text-purple-600" />
              </div>
              <h3 className="font-semibold mb-1 sm:mb-2 text-sm sm:text-base">Beautiful CLI</h3>
              <p className="text-gray-600 text-xs sm:text-sm">
                Colors, spinners, and clear output.
              </p>
            </div>
            <div className="bg-white p-4 sm:p-6 rounded-xl border shadow-sm text-center">
              <div className="h-10 w-10 sm:h-12 sm:w-12 bg-orange-100 rounded-lg flex items-center justify-center mb-3 sm:mb-4 mx-auto">
                <Server className="h-5 w-5 sm:h-6 sm:w-6 text-orange-600" />
              </div>
              <h3 className="font-semibold mb-1 sm:mb-2 text-sm sm:text-base">Any Framework</h3>
              <p className="text-gray-600 text-xs sm:text-sm">
                Node, Python, Go, Ruby - anything.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Advanced Usage */}
      <section id="advanced" className="py-12 sm:py-20 px-4 sm:px-6 lg:px-8 bg-gray-50">
        <div className="max-w-5xl mx-auto">
          <div className="text-center mb-8 sm:mb-12">
            <div className="inline-flex items-center gap-2 bg-gray-200 text-gray-700 px-3 py-1 rounded-full text-xs sm:text-sm font-medium mb-4">
              <Settings className="h-4 w-4" />
              Advanced
            </div>
            <h2 className="text-2xl sm:text-3xl font-bold mb-4">Manual Configuration</h2>
            <p className="text-gray-600 max-w-2xl mx-auto text-sm sm:text-base">
              Need more control? Generate certificates and configure your server.
            </p>
          </div>

          <div className="bg-white rounded-xl border shadow-sm overflow-hidden">
            <div className="p-4 sm:p-6">
              {/* Generate cert command */}
              <div className="mb-6">
                <p className="text-xs sm:text-sm text-gray-600 mb-2">Generate a certificate:</p>
                <div className="bg-gray-900 text-green-400 p-3 sm:p-4 rounded-lg font-mono text-xs sm:text-sm">
                  <code className="break-all">$ instanttls cert myapp.local</code>
                </div>
              </div>

              {/* Code examples - stacked on mobile */}
              <div className="space-y-4">
                {/* Node.js */}
                <div className="bg-gray-900 rounded-lg p-3 sm:p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <Code className="h-4 w-4 text-green-500" />
                    <span className="text-gray-400 text-xs font-medium">Node.js</span>
                  </div>
                  <pre className="text-green-400 font-mono text-xs whitespace-pre-wrap break-all">{`const https = require('https');
const fs = require('fs');

https.createServer({
  key: fs.readFileSync(
    '~/.instanttls/certs/myapp.local/key.pem'
  ),
  cert: fs.readFileSync(
    '~/.instanttls/certs/myapp.local/cert.pem'
  )
}, app).listen(443);`}</pre>
                </div>

                {/* Nginx */}
                <div className="bg-gray-900 rounded-lg p-3 sm:p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <Server className="h-4 w-4 text-blue-500" />
                    <span className="text-gray-400 text-xs font-medium">Nginx</span>
                  </div>
                  <pre className="text-blue-400 font-mono text-xs whitespace-pre-wrap break-all">{`server {
  listen 443 ssl;
  server_name myapp.local;
  ssl_certificate ~/.instanttls/certs/myapp.local/cert.pem;
  ssl_certificate_key ~/.instanttls/certs/myapp.local/key.pem;
  location / {
    proxy_pass http://localhost:3000;
  }
}`}</pre>
                </div>

                {/* Caddy */}
                <div className="bg-gray-900 rounded-lg p-3 sm:p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <Zap className="h-4 w-4 text-purple-500" />
                    <span className="text-gray-400 text-xs font-medium">Caddy</span>
                  </div>
                  <pre className="text-purple-400 font-mono text-xs whitespace-pre-wrap break-all">{`myapp.local {
  tls ~/.instanttls/certs/myapp.local/cert.pem
      ~/.instanttls/certs/myapp.local/key.pem
  reverse_proxy localhost:3000
}`}</pre>
                </div>
              </div>

              <div className="mt-6 p-3 sm:p-4 bg-amber-50 border border-amber-200 rounded-lg">
                <p className="text-amber-800 text-xs sm:text-sm">
                  <strong>💡 Tip:</strong> The <code className="bg-amber-100 px-1 rounded">instanttls serve</code> command 
                  does all this automatically. Only use manual configuration if you need custom server settings.
                </p>
              </div>
            </div>
          </div>

          {/* Certificate paths */}
          <div className="mt-6 sm:mt-8 bg-white rounded-xl border shadow-sm p-4 sm:p-6">
            <h3 className="font-semibold mb-4 text-sm sm:text-base">Certificate Locations</h3>
            <div className="space-y-3 sm:space-y-0 sm:grid sm:grid-cols-2 sm:gap-4">
              <div className="bg-gray-100 p-3 sm:p-4 rounded-lg">
                <p className="text-gray-500 text-xs mb-1">Certificate</p>
                <p className="text-gray-900 font-mono text-xs break-all">~/.instanttls/certs/[domain]/cert.pem</p>
              </div>
              <div className="bg-gray-100 p-3 sm:p-4 rounded-lg">
                <p className="text-gray-500 text-xs mb-1">Private Key</p>
                <p className="text-gray-900 font-mono text-xs break-all">~/.instanttls/certs/[domain]/key.pem</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Pricing */}
      <section className="py-12 sm:py-20 px-4 sm:px-6 lg:px-8 bg-white">
        <div className="max-w-5xl mx-auto">
          <h2 className="text-2xl sm:text-3xl font-bold text-center mb-4">Simple Pricing</h2>
          <p className="text-center text-gray-600 mb-8 sm:mb-12 text-sm sm:text-base">Start free, upgrade when you need more.</p>
          
          <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6 sm:gap-8">
            {/* Free */}
            <div className="bg-white p-6 sm:p-8 rounded-xl border shadow-sm">
              <h3 className="text-lg font-semibold mb-2">Free</h3>
              <div className="text-2xl sm:text-3xl font-bold mb-4">$0<span className="text-base sm:text-lg font-normal text-gray-500">/mo</span></div>
              <ul className="space-y-3 mb-6 sm:mb-8 text-sm sm:text-base">
                <li className="flex items-center gap-2">
                  <CheckCircle className="h-4 w-4 sm:h-5 sm:w-5 text-green-500 flex-shrink-0" />
                  <span>1 wildcard certificate</span>
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle className="h-4 w-4 sm:h-5 sm:w-5 text-green-500 flex-shrink-0" />
                  <span>Local CA generation</span>
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle className="h-4 w-4 sm:h-5 sm:w-5 text-green-500 flex-shrink-0" />
                  <span>Auto-renew</span>
                </li>
              </ul>
              <Link href="/register">
                <Button variant="outline" className="w-full">Get Started</Button>
              </Link>
            </div>

            {/* Pro */}
            <div className="bg-white p-6 sm:p-8 rounded-xl border-2 border-primary shadow-lg relative sm:col-span-2 lg:col-span-1">
              <div className="absolute -top-3 left-1/2 -translate-x-1/2 bg-primary text-white px-3 py-1 rounded-full text-xs sm:text-sm font-medium">
                Popular
              </div>
              <h3 className="text-lg font-semibold mb-2">Pro</h3>
              <div className="text-2xl sm:text-3xl font-bold mb-4">$9<span className="text-base sm:text-lg font-normal text-gray-500">/mo</span></div>
              <ul className="space-y-3 mb-6 sm:mb-8 text-sm sm:text-base">
                <li className="flex items-center gap-2">
                  <CheckCircle className="h-4 w-4 sm:h-5 sm:w-5 text-green-500 flex-shrink-0" />
                  <span>Unlimited certificates</span>
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle className="h-4 w-4 sm:h-5 sm:w-5 text-green-500 flex-shrink-0" />
                  <span>Priority support</span>
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle className="h-4 w-4 sm:h-5 sm:w-5 text-green-500 flex-shrink-0" />
                  <span>Multiple machines</span>
                </li>
              </ul>
              <Link href="/register">
                <Button className="w-full">Upgrade to Pro</Button>
              </Link>
            </div>

            {/* Team */}
            <div className="bg-white p-6 sm:p-8 rounded-xl border shadow-sm">
              <h3 className="text-lg font-semibold mb-2">Team</h3>
              <div className="text-2xl sm:text-3xl font-bold mb-4">$29<span className="text-base sm:text-lg font-normal text-gray-500">/mo</span></div>
              <ul className="space-y-3 mb-6 sm:mb-8 text-sm sm:text-base">
                <li className="flex items-center gap-2">
                  <CheckCircle className="h-4 w-4 sm:h-5 sm:w-5 text-green-500 flex-shrink-0" />
                  <span>Everything in Pro</span>
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle className="h-4 w-4 sm:h-5 sm:w-5 text-green-500 flex-shrink-0" />
                  <span>Team management</span>
                </li>
                <li className="flex items-center gap-2">
                  <CheckCircle className="h-4 w-4 sm:h-5 sm:w-5 text-green-500 flex-shrink-0" />
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
      <footer className="py-8 sm:py-12 px-4 sm:px-6 lg:px-8 border-t">
        <div className="max-w-5xl mx-auto flex flex-col sm:flex-row justify-between items-center gap-4">
          <div className="flex items-center gap-2">
            <Shield className="h-6 w-6 text-primary" />
            <span className="font-bold">InstantTLS</span>
          </div>
          <p className="text-gray-500 text-xs sm:text-sm text-center">
            © 2026 InstantTLS. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  )
}
