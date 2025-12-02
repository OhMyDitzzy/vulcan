import React, { useState } from 'react';
import { getTransaction, getMempool } from '../api/client';
import { Transaction } from '../types';
import { formatAddress } from '../crypto/wallet';

const TransactionExplorer: React.FC = () => {
  const [txId, setTxId] = useState('');
  const [transaction, setTransaction] = useState<Transaction | null>(null);
  const [status, setStatus] = useState<string>('');
  const [error, setError] = useState<string>('');

  const searchTransaction = async () => {
    if (!txId) return;
    try {
      setError('');
      const result = await getTransaction(txId);
      setTransaction(result.transaction);
      setStatus(result.status);
    } catch (err) {
      setError('Transaction not found');
      setTransaction(null);
    }
  };

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">Transaction Explorer</h2>
      
      <div className="card">
        <h3 className="text-xl font-bold mb-4">Search Transaction</h3>
        <div className="flex space-x-2">
          <input
            type="text"
            value={txId}
            onChange={(e) => setTxId(e.target.value)}
            placeholder="Enter transaction ID"
            className="input flex-1"
          />
          <button onClick={searchTransaction} className="btn-primary">
            Search
          </button>
        </div>
        {error && <p className="text-red-600 mt-2">{error}</p>}
      </div>

      {transaction && (
        <div className="card">
          <h3 className="text-xl font-bold mb-4">Transaction Details</h3>
          <div className="space-y-3">
            <div><span className="font-bold">ID:</span> <code className="text-xs">{transaction.id}</code></div>
            <div><span className="font-bold">Status:</span> <span className={status === 'confirmed' ? 'text-green-600' : 'text-orange-600'}>{status}</span></div>
            <div><span className="font-bold">From:</span> <code className="text-xs">{formatAddress(transaction.from)}</code></div>
            <div><span className="font-bold">To:</span> <code className="text-xs">{formatAddress(transaction.to)}</code></div>
            <div><span className="font-bold">Amount:</span> {transaction.amount}</div>
            <div><span className="font-bold">Fee:</span> {transaction.fee}</div>
            <div><span className="font-bold">Timestamp:</span> {new Date(transaction.timestamp).toLocaleString()}</div>
          </div>
        </div>
      )}
    </div>
  );
};

export default TransactionExplorer;