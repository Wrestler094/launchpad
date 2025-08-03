'use client'

import { useState, useEffect, useCallback } from 'react'
import { apiClient } from '@/lib/api'
import { Token, Presale } from '@/types'

export default function Dashboard() {
  const [tokens, setTokens] = useState<Token[]>([])
  const [presales, setPresales] = useState<Presale[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const loadData = useCallback(async () => {
    try {
      setLoading(true)
      const [tokensResponse, presalesResponse] = await Promise.all([
        apiClient.listTokens(),
        apiClient.listPresales()
      ])
      
      setTokens(tokensResponse.data || [])
      setPresales(presalesResponse.data || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load data')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    loadData()
  }, [loadData])

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-indigo-500"></div>
      </div>
    )
  }

  if (error) {
    return (
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
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg leading-6 font-medium text-gray-900">Dashboard</h3>
        <p className="mt-1 text-sm text-gray-500">
          Manage your tokens and presales
        </p>
      </div>

      {/* Tokens Section */}
      <div className="bg-white shadow overflow-hidden sm:rounded-md">
        <div className="px-4 py-5 sm:px-6">
          <h3 className="text-lg leading-6 font-medium text-gray-900">Your Tokens</h3>
          <p className="mt-1 max-w-2xl text-sm text-gray-500">
            Tokens you have created
          </p>
        </div>
        <ul className="divide-y divide-gray-200">
          {tokens.length === 0 ? (
            <li className="px-4 py-4 text-center text-gray-500">
              No tokens created yet
            </li>
          ) : (
            tokens.map((token) => (
              <li key={token.id} className="px-4 py-4">
                <div className="flex items-center justify-between">
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-gray-900 truncate">
                      {token.name} ({token.symbol})
                    </p>
                    <p className="text-sm text-gray-500">
                      Supply: {token.total_supply} | Address: {token.address}
                    </p>
                  </div>
                  <div className="text-sm text-gray-500">
                    {new Date(token.created_at).toLocaleDateString()}
                  </div>
                </div>
              </li>
            ))
          )}
        </ul>
      </div>

      {/* Presales Section */}
      <div className="bg-white shadow overflow-hidden sm:rounded-md">
        <div className="px-4 py-5 sm:px-6">
          <h3 className="text-lg leading-6 font-medium text-gray-900">Your Presales</h3>
          <p className="mt-1 max-w-2xl text-sm text-gray-500">
            Presales you have created
          </p>
        </div>
        <ul className="divide-y divide-gray-200">
          {presales.length === 0 ? (
            <li className="px-4 py-4 text-center text-gray-500">
              No presales created yet
            </li>
          ) : (
            presales.map((presale) => (
              <li key={presale.id} className="px-4 py-4">
                <div className="flex items-center justify-between">
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-gray-900 truncate">
                      Token: {presale.token_address}
                    </p>
                    <p className="text-sm text-gray-500">
                      Rate: {presale.rate} | Soft Cap: {presale.soft_cap} | Hard Cap: {presale.hard_cap}
                    </p>
                    <p className="text-sm text-gray-500">
                      Deadline: {new Date(presale.deadline).toLocaleDateString()}
                    </p>
                  </div>
                  <div className="flex flex-col items-end">
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                      presale.active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
                    }`}>
                      {presale.active ? 'Active' : 'Inactive'}
                    </span>
                    <div className="mt-1 text-sm text-gray-500">
                      {new Date(presale.created_at).toLocaleDateString()}
                    </div>
                  </div>
                </div>
              </li>
            ))
          )}
        </ul>
      </div>
    </div>
  )
}