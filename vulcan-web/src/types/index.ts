// Type definitions matching backend structures

export interface Transaction {
  id: string;
  from: string;
  to: string;
  amount: number;
  fee: number;
  signature: string;
  timestamp: string;
}

export interface Block {
  index: number;
  timestamp: string;
  transactions: Transaction[];
  nonce: number;
  previous_hash: string;
  merkle_root: string;
  hash: string;
  difficulty: number;
}

export interface UTXO {
  tx_id: string;
  address: string;
  amount: number;
  index: number;
}

export interface Wallet {
  address: string;
  private_key: string;
}

export interface BalanceResponse {
  address: string;
  balance: number;
  utxos: UTXO[];
}

export interface HealthResponse {
  status: string;
  height: number;
  mempool: number;
  peers: number;
}

export interface TransactionPayload {
  from: string;
  to: string;
  amount: number;
  fee: number;
}

export interface SignTransactionRequest {
  private_key: string;
  transaction: TransactionPayload;
}

export interface MineRequest {
  miner_address: string;
}

export interface AddPeerRequest {
  address: string;
}

export interface BlocksResponse {
  blocks: Block[];
  start: number;
  limit: number;
  total: number;
}

export interface MempoolResponse {
  transactions: Transaction[];
  count: number;
}

export interface PeersResponse {
  peers: string[];
  count: number;
}