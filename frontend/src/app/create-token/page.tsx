export default function CreateToken() {
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

        <div className="max-w-2xl mx-auto">
          <div className="bg-white/10 backdrop-blur-sm rounded-xl p-8 border border-white/20">
            <h2 className="text-3xl font-bold mb-6 text-center">Create Your Token</h2>
            
            <form className="space-y-6">
              <div>
                <label className="block text-sm font-medium mb-2">Token Name</label>
                <input
                  type="text"
                  placeholder="e.g., My Awesome Token"
                  className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-white placeholder-gray-400"
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">Token Symbol</label>
                <input
                  type="text"
                  placeholder="e.g., MAT"
                  className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-white placeholder-gray-400"
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">Total Supply</label>
                <input
                  type="number"
                  placeholder="e.g., 1000000"
                  className="w-full px-4 py-3 bg-white/10 border border-white/20 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-white placeholder-gray-400"
                />
                <p className="text-sm text-gray-400 mt-1">Total number of tokens to create</p>
              </div>

              <div className="bg-blue-500/20 border border-blue-500/30 rounded-lg p-4">
                <h3 className="font-semibold mb-2">⚡ Quick Info</h3>
                <ul className="text-sm text-gray-300 space-y-1">
                  <li>• Token will be deployed to local Hardhat network</li>
                  <li>• Uses OpenZeppelin ERC-20 standard</li>
                  <li>• You will be the owner of the token</li>
                  <li>• Gas fees will be minimal on local network</li>
                </ul>
              </div>

              <button
                type="submit"
                className="w-full bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 px-8 py-4 rounded-lg font-semibold text-lg transition-all transform hover:scale-[1.02]"
              >
                Deploy Token
              </button>
            </form>
          </div>
        </div>
      </div>
    </main>
  )
}