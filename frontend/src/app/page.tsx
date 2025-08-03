'use client'

import { useState, useEffect } from 'react'
import { useAccount, useConnect, useDisconnect, useSignMessage } from 'wagmi'
import { apiClient } from '@/lib/api'
import TokenCreator from '@/components/TokenCreator'
import PresaleCreator from '@/components/PresaleCreator'
import Dashboard from '@/components/Dashboard'

export default function Home() {
  const { address, isConnected } = useAccount()
  const { connect, connectors } = useConnect()
  const { disconnect } = useDisconnect()
  const { signMessage } = useSignMessage()
  
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [activeTab, setActiveTab] = useState<'dashboard' | 'create-token' | 'create-presale'>('dashboard')

  useEffect(() => {
    // Check if user is already authenticated
    const token = localStorage.getItem('auth_token')
    if (token && isConnected) {
      apiClient.verifyToken()
        .then(() => setIsAuthenticated(true))
        .catch(() => {
          localStorage.removeItem('auth_token')
          apiClient.clearToken()
        })
    }
  }, [isConnected])

  const handleLogin = async () => {
    if (!address) return
    
    setIsLoading(true)
    setError(null)

    try {
      // Get nonce
      const nonceResponse = await apiClient.generateNonce(address)
      const nonce = nonceResponse.data.nonce

      // Sign message
      const message = `Sign this message to authenticate with Launchpad.\n\nNonce: ${nonce}`
      
      signMessage({ message }, {
        onSuccess: async (signature) => {
          try {
            // Login with signature
            const loginResponse = await apiClient.login(address, signature, nonce)
            
            // Set token
            apiClient.setToken(loginResponse.data.token)
            setIsAuthenticated(true)
          } catch (err) {
            setError(err instanceof Error ? err.message : 'Login failed')
          } finally {
            setIsLoading(false)
          }
        },
        onError: () => {
          setError('Signature rejected')
          setIsLoading(false)
        }
      })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to generate nonce')
      setIsLoading(false)
    }
  }

  const handleLogout = () => {
    apiClient.clearToken()
    setIsAuthenticated(false)
    disconnect()
  }

  if (!isConnected) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="max-w-md w-full space-y-8">
          <div className="text-center">
            <h2 className="mt-6 text-3xl font-extrabold text-gray-900">
              Welcome to Launchpad
            </h2>
            <p className="mt-2 text-sm text-gray-600">
              Connect your wallet to get started
            </p>
          </div>
          <div className="space-y-4">
            {connectors.map((connector) => (
              <button
                key={connector.uid}
                onClick={() => connect({ connector })}
                className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                Connect {connector.name}
              </button>
            ))}
          </div>
        </div>
      </div>
    )
  }

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="max-w-md w-full space-y-8">
          <div className="text-center">
            <h2 className="mt-6 text-3xl font-extrabold text-gray-900">
              Authenticate
            </h2>
            <p className="mt-2 text-sm text-gray-600">
              Connected as: {address}
            </p>
            {error && (
              <p className="mt-2 text-sm text-red-600">
                {error}
              </p>
            )}
          </div>
          <div className="space-y-4">
            <button
              onClick={handleLogin}
              disabled={isLoading}
              className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
            >
              {isLoading ? 'Signing...' : 'Sign Message to Login'}
            </button>
            <button
              onClick={handleLogout}
              className="group relative w-full flex justify-center py-2 px-4 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              Disconnect Wallet
            </button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex">
              <div className="flex-shrink-0 flex items-center">
                <h1 className="text-xl font-bold text-gray-900">Launchpad</h1>
              </div>
              <div className="ml-10 flex items-baseline space-x-4">
                <button
                  onClick={() => setActiveTab('dashboard')}
                  className={`px-3 py-2 rounded-md text-sm font-medium ${
                    activeTab === 'dashboard'
                      ? 'bg-indigo-100 text-indigo-700'
                      : 'text-gray-500 hover:text-gray-700'
                  }`}
                >
                  Dashboard
                </button>
                <button
                  onClick={() => setActiveTab('create-token')}
                  className={`px-3 py-2 rounded-md text-sm font-medium ${
                    activeTab === 'create-token'
                      ? 'bg-indigo-100 text-indigo-700'
                      : 'text-gray-500 hover:text-gray-700'
                  }`}
                >
                  Create Token
                </button>
                <button
                  onClick={() => setActiveTab('create-presale')}
                  className={`px-3 py-2 rounded-md text-sm font-medium ${
                    activeTab === 'create-presale'
                      ? 'bg-indigo-100 text-indigo-700'
                      : 'text-gray-500 hover:text-gray-700'
                  }`}
                >
                  Create Presale
                </button>
              </div>
            </div>
            <div className="flex items-center">
              <span className="text-sm text-gray-500 mr-4">
                {address?.slice(0, 6)}...{address?.slice(-4)}
              </span>
              <button
                onClick={handleLogout}
                className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md text-sm font-medium"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          {activeTab === 'dashboard' && <Dashboard />}
          {activeTab === 'create-token' && <TokenCreator />}
          {activeTab === 'create-presale' && <PresaleCreator />}
        </div>
      </main>
    </div>
  )
}
