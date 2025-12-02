import axios from 'axios';
import {
  Block,
  Transaction,
  Wallet,
  BalanceResponse,
  HealthResponse,
  BlocksResponse,
  MempoolResponse,
  PeersResponse,
  SignTransactionRequest,
  MineRequest,
  AddPeerRequest,
  TransactionPayload
} from '../types';

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const getHealth = async (): Promise<HealthResponse> => {
  const response = await api.get('/health');
  return response.data;
};

export const getBlocks = async (start: number = 0, limit: number = 10): Promise<BlocksResponse> => {
  const response = await api.get('/blockchain/blocks', {
    params: { start, limit }
  });
  return response.data;
};

export const getBlock = async (hash: string): Promise<Block> => {
  const response = await api.get(`/blockchain/block/${hash}`);
  return response.data;
};

export const getTransaction = async (txid: string): Promise<{ transaction: Transaction; status: string; block?: string; block_index?: number }> => {
  const response = await api.get(`/blockchain/tx/${txid}`);
  return response.data;
};

export const createWallet = async (): Promise<Wallet> => {
  const response = await api.get('/wallet/new', {
    params: { consent: 'true' }
  });
  return response.data;
};

export const signTransaction = async (request: SignTransactionRequest): Promise<Transaction> => {
  const response = await api.post('/wallet/sign', request);
  return response.data;
};

export const broadcastTransaction = async (transaction: Transaction): Promise<{ message: string; tx_id: string }> => {
  const response = await api.post('/tx', transaction);
  return response.data;
};

export const getMempool = async (): Promise<MempoolResponse> => {
  const response = await api.get('/mempool');
  return response.data;
};

export const mineBlock = async (minerAddress: string): Promise<{ message: string; block: Block }> => {
  const request: MineRequest = { miner_address: minerAddress };
  const response = await api.post('/mine', request);
  return response.data;
};

export const getBalance = async (address: string): Promise<BalanceResponse> => {
  const response = await api.get(`/balance/${address}`);
  return response.data;
};

export const getPeers = async (): Promise<PeersResponse> => {
  const response = await api.get('/peers');
  return response.data;
};

export const addPeer = async (address: string): Promise<{ message: string; address: string }> => {
  const request: AddPeerRequest = { address };
  const response = await api.post('/peers', request);
  return response.data;
};

export const getMetrics = async (): Promise<string> => {
  const response = await api.get('/metrics');
  return response.data;
};

export default api;