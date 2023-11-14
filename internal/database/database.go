package database

import "sync"

type DatabaseInterface interface {
	Get(key string) FileStore
	Set(key string, value FileStore)
	Delete(key string)
	GetAll() map[string]FileStore
}

type FileStore struct {
	Timestamp           string
	ContentType         string
	ContentDisposition  string
}

type Database struct {
	data map[string]FileStore
	mu   sync.RWMutex
}

func NewDatabase() *Database {
	return &Database{
		data: make(map[string]FileStore),
	}
}

func (db *Database) Get(key string) FileStore {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.data[key]
}

func (db *Database) Set(key string, value FileStore) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
}

func (db *Database) Delete(key string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.data, key)
}

func (db *Database) GetAll() map[string]FileStore {
	return db.data
}