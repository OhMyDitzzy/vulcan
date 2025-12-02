import React, { useState } from 'react';
import { createWallet, getBalance, signTransaction as signTxAPI, broadcastTransaction } from '../api/client';
import { generateKeyPair, signTransaction, isValidAddress, formatAddress } from '../crypto/wallet';
import { Wallet, TransactionPayload } from '../types';

const WalletManager: React.FC = () => {
  const [wallet, setWallet] = useState<Wallet | null>(() => {
    const saved = localStorage.getItem('vulcan_wallet');
    return saved ? JSON.parse(saved) : null;
  });
  
  const [balance, setBalance] = useState<number>(0);
  const [showPrivateKey, setShowPrivateKey] = useState(false);
  
  const [txForm, setTxForm] = useState<TransactionPayload>({
    from: wallet?.address || '',
    to: '',
    amount: 0,
    fee: 1
  });

  const handleCreateWallet = () => {
    const newWallet = generateKeyPair();
    const walletData = { address: newWallet.address, private_key: newWallet.privateKey };
 
    setWallet(walletData);
    localStorage.setItem('vulcan_wallet', JSON.stringify(walletData));
    setTxForm({ ...txForm, from: newWallet.address });
  };

  const handleCheckBalance = async () => {
    if (!wallet) return;
    try {
      const data = await getBalance(wallet.address);
      setBalance(data.balance);
    } catch (error) {
      console.error('Failed to fetch balance:', error);
    }
  };

  const handleSendTransaction = async () => {
    if (!wallet) return;
    try {
      const signedTx = signTransaction(txForm, wallet.private_key);
    
      await broadcastTransaction(signedTx);
      alert('Transaction broadcast successfully!');

      setTxForm({ ...txForm, to: '', amount: 0 });
    } catch (error) {
      alert('Failed to send transaction: ' + error);
    }
  };

  const handleDeleteWallet = () => {
    if (confirm('Are you sure you want to delete this wallet? This cannot be undone!')) {
      localStorage.removeItem('vulcan_wallet');
      setWallet(null);
      setBalance(0);
      setTxForm({ from: '', to: '', amount: 0, fee: 1 });
    }
  };

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">Wallet Manager</h2>

      {!wallet ? (
        <div className="card">
          <h3 className="text-xl font-bold mb-4">Create New Wallet</h3>
          <p className="text-gray-600 mb-4">Generate a new wallet with a private key and address.</p>
          <button onClick={handleCreateWallet} className="btn-primary">
            Generate Wallet
          </button>
        </div>
      ) : (
        <>
          <div className="card bg-gradient-to-r from-primary-50 to-primary-100">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-xl font-bold">Your Wallet</h3>
              <button onClick={handleDeleteWallet} className="text-red-600 hover:text-red-800 text-sm">
                🗑️ Delete Wallet
              </button>
            </div>
            <div className="space-y-3">
              <div>
                <label className="label">Address</label>
                <code className="block p-2 bg-white rounded text-xs break-all">{wallet.address}</code>
              </div>
              <div>
                <label className="label">Private Key</label>
                <div className="flex space-x-2">
                  <code className="block p-2 bg-white rounded text-xs break-all flex-1">
                    {showPrivateKey ? wallet.private_key : '••••••••••••••••'}
                  </code>
                  <button onClick={() => setShowPrivateKey(!showPrivateKey)} className="btn-secondary">
                    {showPrivateKey ? '🙈 Hide' : '👁️ Show'}
                  </button>
                </div>
                <p className="text-red-600 text-xs mt-1">⚠️ Never share your private key!</p>
              </div>
              <div className="flex items-center justify-between pt-4 border-t">
                <div>
                  <span className="text-gray-600">Balance:</span>
                  <span className="text-2xl font-bold ml-2">{balance} Vulcan</span>
                </div>
                <button onClick={handleCheckBalance} className="btn-primary">
                  Refresh Balance
                </button>
              </div>
            </div>
          </div>

          <div className="card">
            <h3 className="text-xl font-bold mb-4">Send Transaction</h3>
            <div className="space-y-4">
              <div>
                <label className="label">To Address</label>
                <input
                  type="text"
                  value={txForm.to}
                  onChange={(e) => setTxForm({ ...txForm, to: e.target.value })}
                  className="input"
                  placeholder="Recipient address"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="label">Amount</label>
                  <input
                    type="number"
                    value={txForm.amount}
                    onChange={(e) => setTxForm({ ...txForm, amount: Number(e.target.value) })}
                    className="input"
                  />
                </div>
                <div>
                  <label className="label">Fee</label>
                  <input
                    type="number"
                    value={txForm.fee}
                    onChange={(e) => setTxForm({ ...txForm, fee: Number(e.target.value) })}
                    className="input"
                  />
                </div>
              </div>
              <button
                onClick={handleSendTransaction}
                disabled={!txForm.to || txForm.amount <= 0}
                className="btn-primary w-full disabled:opacity-50"
              >
                Send Transaction
              </button>
            </div>
          </div>
        </>
      )}
    </div>
  );
};

export default WalletManager;