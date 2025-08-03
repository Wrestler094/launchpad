'use client'

import { useState, useEffect, useCallback } from 'react'
import { apiClient } from '@/lib/api'
import { Token } from '@/types'

export default function PresaleCreator() {
  const [tokens, setTokens] = useState<Token[]>([])
  const [formData, setFormData] = useState({
    tokenAddress: '',
    rate: '',
    softCap: '',
    hardCap: '',
    deadline: ''
  })
  const [loading, setLoading] = useState(false)
  const [loadingTokens, setLoadingTokens] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<{ message: string; landingUrl: string } | null>(null)

  const loadTokens = useCallback(async () => {
    try {
      const response = await apiClient.listTokens()
      setTokens(response.data || [])
    } catch {
      setError('Failed to load tokens')
    } finally {
      setLoadingTokens(false)
    }
  }, [])

  useEffect(() => {
    loadTokens()
  }, [loadTokens])

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError(null)
    setSuccess(null)

    try {
      // Validate form
      if (!formData.tokenAddress || !formData.rate || !formData.softCap || !formData.hardCap || !formData.deadline) {
        throw new Error('All fields are required')
      }

      if (isNaN(Number(formData.rate)) || Number(formData.rate) <= 0) {
        throw new Error('Rate must be a positive number')
      }

      if (isNaN(Number(formData.softCap)) || Number(formData.softCap) <= 0) {
        throw new Error('Soft cap must be a positive number')
      }

      if (isNaN(Number(formData.hardCap)) || Number(formData.hardCap) <= Number(formData.softCap)) {
        throw new Error('Hard cap must be greater than soft cap')
      }

      const deadlineDate = new Date(formData.deadline)
      if (deadlineDate <= new Date()) {
        throw new Error('Deadline must be in the future')
      }

      const response = await apiClient.createPresale(
        formData.tokenAddress,
        formData.rate,
        formData.softCap,
        formData.hardCap,
        deadlineDate.toISOString()
      )

      setSuccess({
        message: `Presale created successfully! Address: ${response.data.address}`,
        landingUrl: response.data.landing_url
      })
      setFormData({
        tokenAddress: '',
        rate: '',
        softCap: '',
        hardCap: '',
        deadline: ''
      })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create presale')
    } finally {
      setLoading(false)
    }
  }

  const getMinDatetime = () => {
    const now = new Date()
    now.setMinutes(now.getMinutes() + 1) // At least 1 minute in the future
    return now.toISOString().slice(0, 16)
  }

  if (loadingTokens) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-indigo-500"></div>
      </div>
    )
  }

  return (
    <div className="max-w-md mx-auto bg-white shadow-lg rounded-lg p-6">
      <div className="mb-6">
        <h3 className="text-lg leading-6 font-medium text-gray-900">Create Presale</h3>
        <p className="mt-1 text-sm text-gray-500">
          Launch a presale for your token
        </p>
      </div>

      {error && (
        <div className="mb-4 bg-red-50 border border-red-200 rounded-md p-4">
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

      {success && (
        <div className="mb-4 bg-green-50 border border-green-200 rounded-md p-4">
          <div className="flex">
            <div className="ml-3">
              <h3 className="text-sm font-medium text-green-800">Success</h3>
              <div className="mt-2 text-sm text-green-700">
                <p>{success.message}</p>
                <p className="mt-2">
                  <strong>Landing Page URL:</strong>
                  <br />
                  <a
                    href={success.landingUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-green-600 hover:text-green-800 underline break-all"
                  >
                    {success.landingUrl}
                  </a>
                </p>
              </div>
            </div>
          </div>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="tokenAddress" className="block text-sm font-medium text-gray-700">
            Token
          </label>
          <select
            name="tokenAddress"
            id="tokenAddress"
            value={formData.tokenAddress}
            onChange={handleInputChange}
            className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm p-2 border"
            required
          >
            <option value="">Select a token</option>
            {tokens.map((token) => (
              <option key={token.id} value={token.address}>
                {token.name} ({token.symbol}) - {token.address}
              </option>
            ))}
          </select>
          {tokens.length === 0 && (
            <p className="mt-1 text-sm text-gray-500">
              No tokens available. Create a token first.
            </p>
          )}
        </div>

        <div>
          <label htmlFor="rate" className="block text-sm font-medium text-gray-700">
            Rate (tokens per ETH)
          </label>
          <input
            type="number"
            name="rate"
            id="rate"
            value={formData.rate}
            onChange={handleInputChange}
            placeholder="e.g., 1000"
            className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm p-2 border"
            required
          />
        </div>

        <div>
          <label htmlFor="softCap" className="block text-sm font-medium text-gray-700">
            Soft Cap (ETH)
          </label>
          <input
            type="number"
            step="0.01"
            name="softCap"
            id="softCap"
            value={formData.softCap}
            onChange={handleInputChange}
            placeholder="e.g., 1"
            className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm p-2 border"
            required
          />
        </div>

        <div>
          <label htmlFor="hardCap" className="block text-sm font-medium text-gray-700">
            Hard Cap (ETH)
          </label>
          <input
            type="number"
            step="0.01"
            name="hardCap"
            id="hardCap"
            value={formData.hardCap}
            onChange={handleInputChange}
            placeholder="e.g., 10"
            className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm p-2 border"
            required
          />
        </div>

        <div>
          <label htmlFor="deadline" className="block text-sm font-medium text-gray-700">
            Deadline
          </label>
          <input
            type="datetime-local"
            name="deadline"
            id="deadline"
            value={formData.deadline}
            onChange={handleInputChange}
            min={getMinDatetime()}
            className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm p-2 border"
            required
          />
        </div>

        <button
          type="submit"
          disabled={loading || tokens.length === 0}
          className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
        >
          {loading ? 'Creating Presale...' : 'Create Presale'}
        </button>
      </form>
    </div>
  )
}