import React, { useState, useEffect } from 'react';
import { mineBlock, getBlocks } from '../api/client';

const MinerControl: React.FC = () => {
  const [minerAddress, setMinerAddress] = useState(() => {
    const wallet = localStorage.getItem('vulcan_wallet');
    return wallet ? JSON.parse(wallet).address : '';
  });
  
  const [mining, setMining] = useState(false);
  const [result, setResult] = useState<string>('');
  const [lastBlockHeight, setLastBlockHeight] = useState(0);

  useEffect(() => {
    const checkBlocks = async () => {
      try {
        const data = await getBlocks(0, 1);
        if (data.blocks.length > 0) {
          setLastBlockHeight(data.blocks[0].index);
        }
      } catch (error) {
        console.error('Failed to fetch blocks:', error);
      }
    };
    
    checkBlocks();
    const interval = setInterval(checkBlocks, 5000);
    return () => clearInterval(interval);
  }, []);

  const handleMine = async () => {
    if (!minerAddress) {
      alert('Please enter a miner address');
      return;
    }
    
    setMining(true);
    setResult('⛏️ Mining in progress... This may take a while.');
    
    const startTime = Date.now();
    
    try {
      const data = await mineBlock(minerAddress);
      const duration = ((Date.now() - startTime) / 1000).toFixed(2);
      
      setResult(`✅ Block ${data.block.index} mined successfully!\n` +
                `Hash: ${data.block.hash}\n` +
                `Time: ${duration}s\n` +
                `Transactions: ${data.block.transactions.length}`);
    
      setLastBlockHeight(data.block.index);
      
    } catch (error: any) {
      const duration = ((Date.now() - startTime) / 1000).toFixed(2);

      try {
        const blocksData = await getBlocks(0, 1);
        if (blocksData.blocks.length > 0 && blocksData.blocks[0].index > lastBlockHeight) {
          setResult(`✅ Block ${blocksData.blocks[0].index} mined successfully!\n` +
                    `⚠️ API returned error but block was added to chain\n` +
                    `Hash: ${blocksData.blocks[0].hash}\n` +
                    `Time: ${duration}s`);
          setLastBlockHeight(blocksData.blocks[0].index);
          return;
        }
      } catch (checkError) {
        console.error('Failed to verify mining:', checkError);
      }

      const errorMsg = error.response?.data?.error || error.message || 'Unknown error';
      setResult(`❌ Mining failed: ${errorMsg}\n` +
                `Duration: ${duration}s\n` +
                `Status: ${error.response?.status || 'Network Error'}`);
      
      console.error('Mining error details:', {
        status: error.response?.status,
        data: error.response?.data,
        message: error.message
      });
    } finally {
      setMining(false);
    }
  };

  const loadWalletAddress = () => {
    const wallet = localStorage.getItem('vulcan_wallet');
    if (wallet) {
      const address = JSON.parse(wallet).address;
      setMinerAddress(address);
      setResult('✅ Wallet address loaded');
    } else {
      setResult('❌ No wallet found. Please create a wallet first.');
    }
  };

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">Miner Control</h2>
      
      <div className="card bg-gradient-to-r from-orange-50 to-yellow-50">
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-lg font-bold">Current Block Height</h3>
            <p className="text-3xl font-bold text-primary-600">{lastBlockHeight}</p>
          </div>
          <div className="text-4xl">⛓️</div>
        </div>
      </div>
      
      <div className="card">
        <h3 className="text-xl font-bold mb-4">Mine New Block</h3>
        <p className="text-gray-600 mb-4">
          Mining will select pending transactions from the mempool and attempt to find a valid nonce.
        </p>
        
        <div className="space-y-4">
          <div>
            <label className="label">Reward Address</label>
            <div className="flex gap-2">
              <input
                type="text"
                value={minerAddress}
                onChange={(e) => setMinerAddress(e.target.value)}
                className="input flex-1"
                placeholder="Enter address to receive mining rewards"
              />
              <button
                onClick={loadWalletAddress}
                className="btn-secondary whitespace-nowrap"
              >
                Use My Wallet
              </button>
            </div>
            <p className="text-xs text-gray-500 mt-1">
              Mining rewards will be sent to this address
            </p>
          </div>
          
          <button
            onClick={handleMine}
            disabled={mining || !minerAddress}
            className="btn-primary w-full disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {mining ? '⛏️ Mining...' : '⛏️ Start Mining'}
          </button>
          
          {result && (
            <div className={`p-4 rounded whitespace-pre-wrap break-words font-mono text-sm ${
                result.includes('✅') 
                ? 'bg-green-50 text-green-800 border border-green-200' 
                : result.includes('⛏️')
                ? 'bg-blue-50 text-blue-800 border border-blue-200'
                : 'bg-red-50 text-red-800 border border-red-200'
            }`}>
               {result}
          </div>
          )}
        </div>
      </div>

      <div className="card bg-blue-50 border border-blue-200">
        <h3 className="text-lg font-bold mb-3 flex items-center gap-2">
          <span>ℹ️</span>
          <span>Mining Information</span>
        </h3>
        <ul className="text-sm space-y-2 text-gray-700">
          <li className="flex items-start gap-2">
            <span className="text-blue-600 font-bold">•</span>
            <span>Mining difficulty determines how many leading zeros the block hash must have</span>
          </li>
          <li className="flex items-start gap-2">
            <span className="text-blue-600 font-bold">•</span>
            <span>Higher difficulty means longer mining times but greater security</span>
          </li>
          <li className="flex items-start gap-2">
            <span className="text-blue-600 font-bold">•</span>
            <span>Mining rewards include the block reward (50 VLC) plus transaction fees</span>
          </li>
          <li className="flex items-start gap-2">
            <span className="text-blue-600 font-bold">•</span>
            <span>Only one miner can successfully mine each block</span>
          </li>
          <li className="flex items-start gap-2">
            <span className="text-blue-600 font-bold">•</span>
            <span>If mining takes too long, try creating some transactions first</span>
          </li>
        </ul>
      </div>

      <div className="card bg-yellow-50 border border-yellow-200">
        <h3 className="text-lg font-bold mb-2 flex items-center gap-2">
          <span>⚠️</span>
          <span>Troubleshooting</span>
        </h3>
        <div className="text-sm space-y-2 text-gray-700">
          <p><strong>Mining fails with 500 error?</strong></p>
          <ul className="ml-4 space-y-1">
            <li>• Check if block was actually mined by refreshing the Blocks tab</li>
            <li>• The backend might have mined the block but failed to send response</li>
            <li>• Check backend logs for detailed error messages</li>
          </ul>
        </div>
      </div>
    </div>
  );
};

export default MinerControl;