import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { searchApi, Transaction, Account, Category } from '../services/api'
import { MagnifyingGlassIcon } from '@heroicons/react/24/outline'
import { Link } from 'react-router-dom'

// Search result type
interface SearchResult {
  transactions: Transaction[]
  accounts: Account[]
  categories: Category[]
}

export default function Search() {
  const [keyword, setKeyword] = useState('')
  const [searchType, setSearchType] = useState('all')

  const { data } = useQuery({
    queryKey: ['search', keyword, searchType],
    queryFn: () => searchApi.search({ q: keyword, type: searchType }),
    enabled: keyword.length > 0,
  })

  const results = (data?.data?.data || {}) as SearchResult
  const transactions = results.transactions || []
  const accounts = results.accounts || []
  const categories = results.categories || []

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">搜索</h2>

      <div className="flex gap-2">
        <div className="flex-1 relative">
          <MagnifyingGlassIcon className="w-5 h-5 absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
          <input
            type="text"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            placeholder="搜索交易、账户、分类..."
            className="input pl-10"
          />
        </div>
        <select
          value={searchType}
          onChange={(e) => setSearchType(e.target.value)}
          className="input w-auto"
        >
          <option value="all">全部</option>
          <option value="transactions">交易</option>
          <option value="accounts">账户</option>
          <option value="categories">分类</option>
        </select>
      </div>

      {keyword && (
        <div className="space-y-6">
          {(searchType === 'all' || searchType === 'transactions') && (
            <div className="card">
              <h3 className="font-semibold mb-4">交易记录 ({transactions.length})</h3>
              <div className="space-y-2">
                {transactions.slice(0, 10).map((tx: Transaction) => (
                  <Link
                    key={tx.id}
                    to={`/transactions?id=${tx.id}`}
                    className="flex items-center justify-between py-2 border-b border-gray-100 last:border-0 hover:bg-gray-50 -mx-2 px-2 rounded"
                  >
                    <div>
                      <p className="font-medium">{tx.category?.name || '未分类'}</p>
                      <p className="text-sm text-gray-500">{tx.bill_date} · {tx.account?.name}</p>
                    </div>
                    <p className={`font-semibold ${tx.type === 'income' ? 'text-green-600' : 'text-red-500'}`}>
                      {tx.type === 'income' ? '+' : '-'}¥{tx.amount?.toFixed(2)}
                    </p>
                  </Link>
                ))}
                {transactions.length === 0 && (
                  <p className="text-gray-500 text-center py-2">无匹配结果</p>
                )}
              </div>
            </div>
          )}

          {(searchType === 'all' || searchType === 'accounts') && (
            <div className="card">
              <h3 className="font-semibold mb-4">账户 ({accounts.length})</h3>
              <div className="space-y-2">
                {accounts.map((acc: Account) => (
                  <Link
                    key={acc.id}
                    to={`/accounts?id=${acc.id}`}
                    className="flex items-center justify-between py-2 border-b border-gray-100 last:border-0 hover:bg-gray-50 -mx-2 px-2 rounded"
                  >
                    <span className="font-medium">{acc.name}</span>
                    <span className="text-gray-500">¥{acc.balance?.toFixed(2)}</span>
                  </Link>
                ))}
                {accounts.length === 0 && (
                  <p className="text-gray-500 text-center py-2">无匹配结果</p>
                )}
              </div>
            </div>
          )}

          {(searchType === 'all' || searchType === 'categories') && (
            <div className="card">
              <h3 className="font-semibold mb-4">分类 ({categories.length})</h3>
              <div className="flex flex-wrap gap-2">
                {categories.map((cat: Category) => (
                  <Link
                    key={cat.id}
                    to={`/categories?id=${cat.id}`}
                    className="flex items-center gap-2 px-3 py-1 bg-gray-100 rounded-full hover:bg-gray-200"
                  >
                    <div
                      className="w-3 h-3 rounded-full"
                      style={{ backgroundColor: cat.color || '#999' }}
                    />
                    <span>{cat.name}</span>
                  </Link>
                ))}
                {categories.length === 0 && (
                  <p className="text-gray-500">无匹配结果</p>
                )}
              </div>
            </div>
          )}
        </div>
      )}

      {!keyword && (
        <div className="text-center py-12 text-gray-500">
          <p>输入关键词开始搜索</p>
        </div>
      )}
    </div>
  )
}
