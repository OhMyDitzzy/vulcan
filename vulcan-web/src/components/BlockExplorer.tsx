import React, { useState, useEffect } from 'react';
import { getBlocks } from '../api/client';
import { Block } from '../types';

const BlockExplorer: React.FC = () => {
  const [blocks, setBlocks] = useState<Block[]>([]);
  const [selectedBlock, setSelectedBlock] = useState<Block | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchBlocks = async () => {
      try {
        const data = await getBlocks(0, 20);
        setBlocks(data.blocks.reverse());
      } catch (error) {
        console.error('Failed to fetch blocks:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchBlocks();
  }, []);

  if (loading) return <div>Loading blocks...</div>;

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">Block Explorer</h2>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card">
          <h3 className="text-xl font-bold mb-4">Blocks</h3>
          <div className="space-y-2 max-h-[600px] overflow-y-auto">
            {blocks.map((block) => (
              <div
                key={block.hash}
                onClick={() => setSelectedBlock(block)}
                className="p-3 border rounded cursor-pointer hover:bg-gray-50"
              >
                <div className="font-bold">Block #{block.index}</div>
                <div className="text-sm text-gray-600 truncate">{block.hash}</div>
                <div className="text-xs text-gray-500">{block.transactions.length} transactions</div>
              </div>
            ))}
          </div>
        </div>
        
        <div className="card">
          <h3 className="text-xl font-bold mb-4">Block Details</h3>
          {selectedBlock ? (
            <div className="space-y-3">
              <div><span className="font-bold">Index:</span> {selectedBlock.index}</div>
              <div><span className="font-bold">Hash:</span> <code className="text-xs">{selectedBlock.hash}</code></div>
              <div><span className="font-bold">Previous Hash:</span> <code className="text-xs">{selectedBlock.previous_hash}</code></div>
              <div><span className="font-bold">Merkle Root:</span> <code className="text-xs">{selectedBlock.merkle_root}</code></div>
              <div><span className="font-bold">Nonce:</span> {selectedBlock.nonce}</div>
              <div><span className="font-bold">Difficulty:</span> {selectedBlock.difficulty}</div>
              <div><span className="font-bold">Timestamp:</span> {new Date(selectedBlock.timestamp).toLocaleString()}</div>
              <div><span className="font-bold">Transactions:</span> {selectedBlock.transactions.length}</div>
            </div>
          ) : (
            <p className="text-gray-500">Select a block to view details</p>
          )}
        </div>
      </div>
    </div>
  );
};

export default BlockExplorer;