'use client'

import { useState, useEffect, useCallback } from 'react'
import { useAccount, useConnect } from 'wagmi'
import { apiClient } from '@/lib/api'
import { PresaleData } from '@/types'

interface PresalePageProps {
  params: Promise<{ id: string }>
}

export default function PresalePage({ params }: PresalePageProps) {
  const { address, isConnected } = useAccount()
  const { connect, connectors } = useConnect()
  
  const [presaleData, setPresaleData] = useState<PresaleData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [purchaseAmount, setPurchaseAmount] = useState('')
  const [purchasing, setPurchasing] = useState(false)
  const [purchaseSuccess, setPurchaseSuccess] = useState<string | null>(null)
  const [presaleId, setPresaleId] = useState<string>('')

  const loadPresaleData = useCallback(async (id: string) => {
    try {
      setLoading(true)
      const response = await apiClient.getPublicPresale(parseInt(id))
      setPresaleData(response.data)
    } catch {
      setError('Failed to load presale data')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    params.then((resolvedParams) => {
      setPresaleId(resolvedParams.id)
      loadPresaleData(resolvedParams.id)
    })
  }, [params, loadPresaleData])

  const handlePurchase = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!isConnected || !address) {
      alert('Please connect your wallet first')
      return
    }

    setPurchasing(true)
    setPurchaseSuccess(null)
    setError(null)

    try {
      if (!purchaseAmount || isNaN(Number(purchaseAmount)) || Number(purchaseAmount) <= 0) {
        throw new Error('Please enter a valid amount')
      }

      // In a real implementation, you would send a transaction to the presale contract
      // For now, we'll simulate this with a mock transaction hash
      const mockTxHash = '0x' + Math.random().toString(16).substr(2, 64)
      
      const response = await apiClient.participateInPresale(
        parseInt(presaleId),
        purchaseAmount,
        mockTxHash
      )

      setPurchaseSuccess(`Purchase successful! You will receive ${response.data.amount_tokens} tokens.`)
      setPurchaseAmount('')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Purchase failed')
    } finally {
      setPurchasing(false)
    }
  }

  const calculateTokens = () => {
    if (!purchaseAmount || !presaleData) return '0'
    const ethAmount = Number(purchaseAmount)
    const rate = Number(presaleData.presale.rate)
    return (ethAmount * rate).toLocaleString()
  }

  const isPresaleActive = () => {
    if (!presaleData) return false
    const now = new Date()
    const deadline = new Date(presaleData.presale.deadline)
    return presaleData.presale.active && now < deadline
  }

  const getTimeRemaining = () => {
    if (!presaleData) return ''
    const now = new Date()
    const deadline = new Date(presaleData.presale.deadline)
    const diff = deadline.getTime() - now.getTime()
    
    if (diff <= 0) return 'Expired'
    
    const days = Math.floor(diff / (1000 * 60 * 60 * 24))
    const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
    const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
    
    return `${days}d ${hours}h ${minutes}m`
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-indigo-500"></div>
      </div>
    )
  }

  if (error && !presaleData) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="max-w-md w-full bg-white shadow-lg rounded-lg p-6">
          <div className="bg-red-50 border border-red-200 rounded-md p-4">
            <div className="flex">
              <div className="ml-3">
                <h3 className="text-sm font-medium text-red-800">Error</h3>
                <div className="mt-2 text-sm text-red-700">
                  <p>{error}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="bg-white shadow-lg rounded-lg overflow-hidden">
          {/* Header */}
          <div className="bg-indigo-600 px-6 py-4">
            <h1 className="text-2xl font-bold text-white">Token Presale</h1>
            <p className="text-indigo-100">
              {presaleData?.token.name} ({presaleData?.token.symbol})
            </p>
          </div>

          {/* Presale Info */}
          <div className="px-6 py-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
              <div className="bg-gray-50 p-4 rounded-lg">
                <h3 className="text-sm font-medium text-gray-500">Token Address</h3>
                <p className="text-sm text-gray-900 break-all">{presaleData?.token.address}</p>
              </div>
              <div className="bg-gray-50 p-4 rounded-lg">
                <h3 className="text-sm font-medium text-gray-500">Total Supply</h3>
                <p className="text-sm text-gray-900">{Number(presaleData?.token.total_supply).toLocaleString()}</p>
              </div>
              <div className="bg-gray-50 p-4 rounded-lg">
                <h3 className="text-sm font-medium text-gray-500">Rate</h3>
                <p className="text-sm text-gray-900">{presaleData?.presale.rate} tokens per ETH</p>
              </div>
              <div className="bg-gray-50 p-4 rounded-lg">
                <h3 className="text-sm font-medium text-gray-500">Time Remaining</h3>
                <p className="text-sm text-gray-900">{getTimeRemaining()}</p>
              </div>
              <div className="bg-gray-50 p-4 rounded-lg">
                <h3 className="text-sm font-medium text-gray-500">Soft Cap</h3>
                <p className="text-sm text-gray-900">{presaleData?.presale.soft_cap} ETH</p>
              </div>
              <div className="bg-gray-50 p-4 rounded-lg">
                <h3 className="text-sm font-medium text-gray-500">Hard Cap</h3>
                <p className="text-sm text-gray-900">{presaleData?.presale.hard_cap} ETH</p>
              </div>
            </div>

            {/* Status */}
            <div className="mb-6">
              <span className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${
                isPresaleActive()
                  ? 'bg-green-100 text-green-800'
                  : 'bg-red-100 text-red-800'
              }`}>
                {isPresaleActive() ? 'Active' : 'Inactive'}
              </span>
            </div>

            {/* Purchase Form */}
            {isPresaleActive() && (
              <div className="border-t pt-6">
                <h3 className="text-lg font-medium text-gray-900 mb-4">Purchase Tokens</h3>
                
                {!isConnected ? (
                  <div className="space-y-4">
                    <p className="text-sm text-gray-600">Connect your wallet to participate</p>
                    {connectors.map((connector) => (
                      <button
                        key={connector.uid}
                        onClick={() => connect({ connector })}
                        className="w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                      >
                        Connect {connector.name}
                      </button>
                    ))}
                  </div>
                ) : (
                  <form onSubmit={handlePurchase} className="space-y-4">
                    <div>
                      <label htmlFor="amount" className="block text-sm font-medium text-gray-700">
                        Amount (ETH)
                      </label>
                      <input
                        type="number"
                        step="0.01"
                        name="amount"
                        id="amount"
                        value={purchaseAmount}
                        onChange={(e) => setPurchaseAmount(e.target.value)}
                        placeholder="0.1"
                        className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm p-2 border"
                        required
                      />
                      {purchaseAmount && (
                        <p className="mt-1 text-sm text-gray-500">
                          You will receive approximately {calculateTokens()} tokens
                        </p>
                      )}
                    </div>

                    <button
                      type="submit"
                      disabled={purchasing}
                      className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                    >
                      {purchasing ? 'Processing...' : 'Buy Tokens'}
                    </button>
                  </form>
                )}
              </div>
            )}

            {/* Messages */}
            {error && (
              <div className="mt-4 bg-red-50 border border-red-200 rounded-md p-4">
                <div className="flex">
                  <div className="ml-3">
                    <h3 className="text-sm font-medium text-red-800">Error</h3>
                    <div className="mt-2 text-sm text-red-700">
                      <p>{error}</p>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {purchaseSuccess && (
              <div className="mt-4 bg-green-50 border border-green-200 rounded-md p-4">
                <div className="flex">
                  <div className="ml-3">
                    <h3 className="text-sm font-medium text-green-800">Success</h3>
                    <div className="mt-2 text-sm text-green-700">
                      <p>{purchaseSuccess}</p>
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}