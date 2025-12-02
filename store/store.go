package store

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
)

// Store provides persistence layer without knowing about domain types
type Store interface {
	SaveBlock(index uint64, hash string, data []byte) error
	GetBlock(index uint64) ([]byte, error)
	GetBlockByHash(hash string) ([]byte, error)
	GetHeight() (uint64, error)
	Close() error
}

type BadgerStore struct {
	db *badger.DB
}

func NewBadgerStore(path string) (*BadgerStore, error) {
	opts := badger.DefaultOptions(path)
	opts.Logger = nil // Suppress logs
	
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	
	return &BadgerStore{db: db}, nil
}

func (bs *BadgerStore) SaveBlock(index uint64, hash string, data []byte) error {
	return bs.db.Update(func(txn *badger.Txn) error {
		indexKey := []byte(fmt.Sprintf("block:index:%d", index))
		if err := txn.Set(indexKey, data); err != nil {
			return err
		}

		hashKey := []byte(fmt.Sprintf("block:hash:%s", hash))
		if err := txn.Set(hashKey, data); err != nil {
			return err
		}

		heightKey := []byte("blockchain:height")
		heightData, _ := json.Marshal(index)
		return txn.Set(heightKey, heightData)
	})
}

func (bs *BadgerStore) GetBlock(index uint64) ([]byte, error) {
	var data []byte
	err := bs.db.View(func(txn *badger.Txn) error {
		key := []byte(fmt.Sprintf("block:index:%d", index))
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		
		data, err = item.ValueCopy(nil)
		return err
	})
	return data, err
}

func (bs *BadgerStore) GetBlockByHash(hash string) ([]byte, error) {
	var data []byte
	err := bs.db.View(func(txn *badger.Txn) error {
		key := []byte(fmt.Sprintf("block:hash:%s", hash))
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		
		data, err = item.ValueCopy(nil)
		return err
	})
	return data, err
}

func (bs *BadgerStore) GetHeight() (uint64, error) {
	var height uint64
	err := bs.db.View(func(txn *badger.Txn) error {
		key := []byte("blockchain:height")
		item, err := txn.Get(key)
		if err == badger.ErrKeyNotFound {
			return nil
		}
		if err != nil {
			return err
		}
		
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &height)
		})
	})
	return height, err
}

func (bs *BadgerStore) Close() error {
	return bs.db.Close()
}