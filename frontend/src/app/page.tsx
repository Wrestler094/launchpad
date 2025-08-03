export default function Home() {
  return (
    <main className="min-h-screen bg-gradient-to-br from-blue-900 via-purple-900 to-indigo-900 text-white">
      <div className="container mx-auto px-4 py-8">
        {/* Header */}
        <header className="flex justify-between items-center mb-12">
          <h1 className="text-4xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-purple-400">
            🚀 Launchpad
          </h1>
          <div className="flex items-center space-x-4">
            <button className="bg-blue-600 hover:bg-blue-700 px-6 py-2 rounded-lg font-semibold transition-colors">
              Connect Wallet
            </button>
          </div>
        </header>

        {/* Hero Section */}
        <div className="text-center mb-16">
          <h2 className="text-6xl font-bold mb-6">
            Launch Your Token
            <br />
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-purple-400 to-pink-400">
              Start Your Presale
            </span>
          </h2>
          <p className="text-xl text-gray-300 mb-8 max-w-3xl mx-auto">
            Create ERC-20 tokens, configure presales, and launch your project with ease. 
            Our platform makes token deployment and presale management simple and secure.
          </p>
          <div className="flex justify-center space-x-4">
            <button className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 px-8 py-4 rounded-lg font-semibold text-lg transition-all transform hover:scale-105">
              Create Token
            </button>
            <button className="border border-purple-400 hover:bg-purple-400 hover:text-purple-900 px-8 py-4 rounded-lg font-semibold text-lg transition-all">
              Learn More
            </button>
          </div>
        </div>

        {/* Features Section */}
        <div className="grid md:grid-cols-3 gap-8 mb-16">
          <div className="bg-white/10 backdrop-blur-sm rounded-xl p-6 border border-white/20">
            <div className="text-3xl mb-4">🪙</div>
            <h3 className="text-xl font-semibold mb-3">Create ERC-20 Tokens</h3>
            <p className="text-gray-300">
              Deploy your own ERC-20 token with custom name, symbol, and supply. 
              Built on proven OpenZeppelin contracts for security.
            </p>
          </div>
          
          <div className="bg-white/10 backdrop-blur-sm rounded-xl p-6 border border-white/20">
            <div className="text-3xl mb-4">🎯</div>
            <h3 className="text-xl font-semibold mb-3">Launch Presales</h3>
            <p className="text-gray-300">
              Configure and launch token presales with customizable rates, caps, 
              and deadlines. Automatic refunds if soft cap isn&apos;t reached.
            </p>
          </div>
          
          <div className="bg-white/10 backdrop-blur-sm rounded-xl p-6 border border-white/20">
            <div className="text-3xl mb-4">🌐</div>
            <h3 className="text-xl font-semibold mb-3">Landing Pages</h3>
            <p className="text-gray-300">
              Get shareable landing pages for your presale automatically. 
              Let investors easily participate with MetaMask integration.
            </p>
          </div>
        </div>

        {/* Quick Start Steps */}
        <div className="bg-white/5 backdrop-blur-sm rounded-xl p-8 border border-white/20">
          <h3 className="text-2xl font-semibold mb-6 text-center">Quick Start</h3>
          <div className="grid md:grid-cols-4 gap-6">
            <div className="text-center">
              <div className="bg-blue-600 w-12 h-12 rounded-full flex items-center justify-center mx-auto mb-3 text-xl font-bold">
                1
              </div>
              <h4 className="font-semibold mb-2">Connect Wallet</h4>
              <p className="text-sm text-gray-300">Connect your MetaMask wallet to get started</p>
            </div>
            
            <div className="text-center">
              <div className="bg-purple-600 w-12 h-12 rounded-full flex items-center justify-center mx-auto mb-3 text-xl font-bold">
                2
              </div>
              <h4 className="font-semibold mb-2">Create Token</h4>
              <p className="text-sm text-gray-300">Define your token parameters and deploy</p>
            </div>
            
            <div className="text-center">
              <div className="bg-pink-600 w-12 h-12 rounded-full flex items-center justify-center mx-auto mb-3 text-xl font-bold">
                3
              </div>
              <h4 className="font-semibold mb-2">Setup Presale</h4>
              <p className="text-sm text-gray-300">Configure presale parameters and launch</p>
            </div>
            
            <div className="text-center">
              <div className="bg-indigo-600 w-12 h-12 rounded-full flex items-center justify-center mx-auto mb-3 text-xl font-bold">
                4
              </div>
              <h4 className="font-semibold mb-2">Share & Sell</h4>
              <p className="text-sm text-gray-300">Share your landing page with investors</p>
            </div>
          </div>
        </div>
      </div>
    </main>
  )
}
