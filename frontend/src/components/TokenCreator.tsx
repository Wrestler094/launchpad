'use client'

import { useState } from 'react'
import { apiClient } from '@/lib/api'

export default function TokenCreator() {
  const [formData, setFormData] = useState({
    name: '',
    symbol: '',
    totalSupply: ''
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
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
      if (!formData.name || !formData.symbol || !formData.totalSupply) {
        throw new Error('All fields are required')
      }

      if (isNaN(Number(formData.totalSupply)) || Number(formData.totalSupply) <= 0) {
        throw new Error('Total supply must be a positive number')
      }

      const response = await apiClient.createToken(
        formData.name,
        formData.symbol,
        formData.totalSupply
      )

      setSuccess(`Token created successfully! Address: ${response.data.address}`)
      setFormData({ name: '', symbol: '', totalSupply: '' })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create token')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-md mx-auto bg-white shadow-lg rounded-lg p-6">
      <div className="mb-6">
        <h3 className="text-lg leading-6 font-medium text-gray-900">Create Token</h3>
        <p className="mt-1 text-sm text-gray-500">
          Deploy a new ERC-20 token
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
                <p>{success}</p>
              </div>
            </div>
          </div>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="name" className="block text-sm font-medium text-gray-700">
            Token Name
          </label>
          <input
            type="text"
            name="name"
            id="name"
            value={formData.name}
            onChange={handleInputChange}
            placeholder="e.g., My Token"
            className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm p-2 border"
            required
          />
        </div>

        <div>
          <label htmlFor="symbol" className="block text-sm font-medium text-gray-700">
            Token Symbol
          </label>
          <input
            type="text"
            name="symbol"
            id="symbol"
            value={formData.symbol}
            onChange={handleInputChange}
            placeholder="e.g., MTK"
            className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm p-2 border"
            required
          />
        </div>

        <div>
          <label htmlFor="totalSupply" className="block text-sm font-medium text-gray-700">
            Total Supply
          </label>
          <input
            type="number"
            name="totalSupply"
            id="totalSupply"
            value={formData.totalSupply}
            onChange={handleInputChange}
            placeholder="e.g., 1000000"
            className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm p-2 border"
            required
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
        >
          {loading ? 'Creating Token...' : 'Create Token'}
        </button>
      </form>
    </div>
  )
}