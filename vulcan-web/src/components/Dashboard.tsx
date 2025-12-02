import React, { useState, useEffect } from 'react';
import { getHealth, getBlocks, getMempool, getPeers } from '../api/client';
import { HealthResponse, Block } from '../types';

const Dashboard: React.FC = () => {
  const [health, setHealth] = useState<HealthResponse | null>(null);
  const [latestBlocks, setLatestBlocks] = useState<Block[]>([]);
  const [mempoolSize, setMempoolSize] = useState(0);
  const [peersCount, setPeersCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchData = async () => {
    try {
      setLoading(true);
      setError(null);

      const healthData = await getHealth();
      setHealth(healthData);

      const blocksData = await getBlocks(0, 5);
      setLatestBlocks(blocksData.blocks.reverse());

      const mempoolData = await getMempool();
      setMempoolSize(mempoolData.count);

      const peersData = await getPeers();
      setPeersCount(peersData.count);
    } catch (err) {
      setError('Failed to fetch dashboard data');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 10000);
    return () => clearInterval(interval);
  }, []);

  if (loading && !health) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="text-lg text-gray-600">Loading dashboard...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4">
        <p className="text-red-800">{error}</p>
        <button onClick={fetchData} className="btn-primary mt-2">
          Retry
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-gray-800">Dashboard</h2>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div className="card bg-gradient-to-br from-primary-500 to-primary-600 text-white">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-primary-100 text-sm">Blockchain Height</p>
              <p className="text-3xl font-bold mt-1">{health?.height}</p>
            </div>
            <div className="text-5xl opacity-20">⛓️</div>
          </div>
        </div>

        <div className="card bg-gradient-to-br from-green-500 to-green-600 text-white">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-green-100 text-sm">Pending Transactions</p>
              <p className="text-3xl font-bold mt-1">{mempoolSize}</p>
            </div>
            <div className="text-5xl opacity-20">💸</div>
          </div>
        </div>

        <div className="card bg-gradient-to-br from-purple-500 to-purple-600 text-white">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-purple-100 text-sm">Connected Peers</p>
              <p className="text-3xl font-bold mt-1">{peersCount}</p>
            </div>
            <div className="text-5xl opacity-20">🌐</div>
          </div>
        </div>

        <div className="card bg-gradient-to-br from-orange-500 to-orange-600 text-white">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-orange-100 text-sm">Node Status</p>
              <p className="text-lg font-bold mt-1">
                {health?.status === 'healthy' ? '🟢 Healthy' : '🔴 Offline'}
              </p>
            </div>
            <div className="text-5xl opacity-20">⚡</div>
          </div>
        </div>
      </div>

      <div className="card">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-xl font-bold text-gray-800">Latest Blocks</h3>
          <button onClick={fetchData} className="text-primary-600 hover:text-primary-700">
            🔄 Refresh
          </button>
        </div>

        {latestBlocks.length === 0 ? (
          <p className="text-gray-500 text-center py-8">No blocks yet</p>
        ) : (
          <div className="space-y-3">
            {latestBlocks.map((block) => (
              <div
                key={block.hash}
                className="border border-gray-200 rounded-lg p-4 hover:bg-gray-50 transition-colors"
              >
                <div className="flex items-center justify-between">
                  <div className="flex-1">
                    <div className="flex items-center space-x-2">
                      <span className="font-mono font-bold text-primary-600">
                        Block #{block.index}
                      </span>
                      <span className="text-xs text-gray-500">
                        {new Date(block.timestamp).toLocaleString()}
                      </span>
                    </div>
                    <div className="mt-1 text-sm font-mono text-gray-600 truncate">
                      Hash: {block.hash}
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-sm text-gray-600">
                      {block.transactions.length} tx
                    </div>
                    <div className="text-xs text-gray-500">
                      Nonce: {block.nonce}
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="card">
        <h3 className="text-xl font-bold text-gray-800 mb-4">Quick Actions</h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <button className="btn-primary py-4">
            <div className="text-2xl mb-2">👛</div>
            <div>Create Wallet</div>
          </button>
          <button className="btn-primary py-4">
            <div className="text-2xl mb-2">💸</div>
            <div>Send Transaction</div>
          </button>
          <button className="btn-primary py-4">
            <div className="text-2xl mb-2">⛏️</div>
            <div>Mine Block</div>
          </button>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;