import * as elliptic from 'elliptic';
import CryptoJS from 'crypto-js';
import { Transaction, TransactionPayload } from '../types';

const ec = new elliptic.ec('secp256k1');

export const generateKeyPair = (): { privateKey: string; address: string } => {
  const keyPair = ec.genKeyPair();
  const privateKey = keyPair.getPrivate('hex');
  const publicKey = keyPair.getPublic('hex');
  
  return {
    privateKey,
    address: publicKey
  };
};

export const signData = (data: string, privateKeyHex: string): string => {
  const keyPair = ec.keyFromPrivate(privateKeyHex, 'hex');
  const dataHash = CryptoJS.SHA256(data).toString();
  const signature = keyPair.sign(dataHash);
  
  return signature.toDER('hex');
};

export const createTransactionData = (tx: TransactionPayload, timestamp: string): string => {
  return `${tx.from}${tx.to}${tx.amount}${tx.fee}${timestamp}`;
};

export const signTransaction = (
  tx: TransactionPayload,
  privateKeyHex: string
): Transaction => {
  const timestamp = new Date().toISOString();
  
  const dataToSign = createTransactionData(tx, timestamp);

  const hash = CryptoJS.SHA256(dataToSign).toString();
  
  const signature = signData(hash, privateKeyHex);
  
  const signedTx: Transaction = {
    id: '',
    from: tx.from,
    to: tx.to,
    amount: tx.amount,
    fee: tx.fee,
    signature: signature,
    timestamp: timestamp
  };
  
  const txData = `${signedTx.from}${signedTx.to}${signedTx.amount}${signedTx.fee}${signedTx.signature}${signedTx.timestamp}`;
  signedTx.id = CryptoJS.SHA256(txData).toString();
  
  return signedTx;
};

export const verifyKeyPair = (privateKeyHex: string, address: string): boolean => {
  try {
    const keyPair = ec.keyFromPrivate(privateKeyHex, 'hex');
    const publicKey = keyPair.getPublic('hex');
    return publicKey === address;
  } catch (error) {
    return false;
  }
};

export const formatAddress = (address: string): string => {
  if (address.length <= 16) return address;
  return `${address.slice(0, 8)}...${address.slice(-8)}`;
};

export const isValidAddress = (address: string): boolean => {
  return /^[0-9a-fA-F]{130}$/.test(address);
};

export const isValidPrivateKey = (privateKey: string): boolean => {
  return /^[0-9a-fA-F]{64}$/.test(privateKey);
};