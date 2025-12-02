import React, { useState } from 'react';
import Dashboard from './components/Dashboard';
import BlockExplorer from './components/BlockExplorer';
import TransactionExplorer from './components/TransactionExplorer';
import WalletManager from './components/WalletManager';
import MinerControl from './components/MinerControl';

type TabType = 'dashboard' | 'blocks' | 'transactions' | 'wallet' | 'miner';

function App() {
  const [activeTab, setActiveTab] = useState<TabType>('dashboard');

  const tabs = [
    { id: 'dashboard' as TabType, label: 'Dashboard', icon: '📊' },
    { id: 'blocks' as TabType, label: 'Blocks', icon: '⛓️' },
    { id: 'transactions' as TabType, label: 'Transactions', icon: '💸' },
    { id: 'wallet' as TabType, label: 'Wallet', icon: '👛' },
    { id: 'miner' as TabType, label: 'Miner', icon: '⛏️' },
  ];

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-gradient-to-r from-primary-600 to-primary-800 text-white shadow-lg">
        <div className="max-w-5xl mx-auto px-4 py-6">
          <div className="flex items-center justify-between flex-wrap gap-4">
            <div>
              <h1 className="text-2xl font-bold">🌋 Vulcan Blockchain</h1>
              <p className="text-primary-100 mt-1 text-sm">
                Decentralized. Transparent. Secure.
              </p>
            </div>
            <div className="text-right">
              <p className="text-sm text-primary-100">Mini Blockchain Explorer</p>
              <p className="text-xs text-primary-200">v1.0.0</p>
            </div>
          </div>
        </div>
      </header>

      <nav className="bg-white shadow-md sticky top-0 z-10">
        <div className="max-w-5xl mx-auto px-2">
          <div className="flex flex-wrap items-center justify-around gap-1 py-2">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`
                  flex flex-col items-center px-3 py-2 text-sm rounded-md transition
                  ${activeTab === tab.id
                    ? 'text-primary-600 font-semibold border-b-2 border-primary-600'
                    : 'text-gray-600 hover:text-gray-900'
                  }
                `}
              >
                <span className="text-lg">{tab.icon}</span>
                {tab.label}
              </button>
            ))}
          </div>
        </div>
      </nav>

      <main className="max-w-5xl mx-auto px-4 py-6">
        {activeTab === 'dashboard' && <Dashboard />}
        {activeTab === 'blocks' && <BlockExplorer />}
        {activeTab === 'transactions' && <TransactionExplorer />}
        {activeTab === 'wallet' && <WalletManager />}
        {activeTab === 'miner' && <MinerControl />}
      </main>

      <footer className="bg-gray-800 text-gray-300 mt-12">
        <div className="max-w-5xl mx-auto px-4 py-6">
          <div className="text-center">
            <p className="text-sm">Built with ❤️ By DitzDev</p>
            <p className="text-xs text-gray-400 mt-2">Powered by Go and React</p>
          </div>
        </div>
      </footer>
    </div>
  );
}

export default App;