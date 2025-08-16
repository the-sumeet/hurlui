package main

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

func (a *App) SetCache(key string, result HurlResult) error {
	if a.cacheDB == nil {
		return fmt.Errorf("cache not initialized")
	}

	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal HurlResult: %w", err)
	}

	return a.cacheDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("hurl_cache"))
		if b == nil {
			return fmt.Errorf("cache bucket does not exist")
		}
		return b.Put([]byte(key), data)
	})
}

func (a *App) GetCache(key string) (HurlResult, error) {
	if a.cacheDB == nil {
		return HurlResult{}, fmt.Errorf("cache not initialized")
	}

	var result HurlResult
	err := a.cacheDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("hurl_cache"))
		if b == nil {
			return fmt.Errorf("cache bucket does not exist")
		}
		data := b.Get([]byte(key))
		if data == nil {
			return fmt.Errorf("key %s not found", key)
		}
		return json.Unmarshal(data, &result)
	})

	return result, err
}

func (a *App) DeleteCache(key string) error {
	if a.cacheDB == nil {
		return fmt.Errorf("cache not initialized")
	}

	return a.cacheDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("hurl_cache"))
		if b == nil {
			return fmt.Errorf("cache bucket does not exist")
		}
		return b.Delete([]byte(key))
	})
}

func (a *App) ClearCache() error {
	if a.cacheDB == nil {
		return fmt.Errorf("cache not initialized")
	}

	return a.cacheDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("hurl_cache"))
		if b == nil {
			return fmt.Errorf("cache bucket does not exist")
		}
		return b.ForEach(func(k, v []byte) error {
			return b.Delete(k)
		})
	})
}

func (a *App) CloseCache() error {
	if a.cacheDB != nil {
		return a.cacheDB.Close()
	}
	return nil
}

// ExistsInCache checks if a key exists in the cache
func (a *App) ExistsInCache(key string) bool {
	if a.cacheDB == nil {
		return false
	}

	exists := false
	a.cacheDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("hurl_cache"))
		if b == nil {
			return nil
		}
		data := b.Get([]byte(key))
		exists = data != nil
		return nil
	})

	return exists
}

// GetCacheKeys returns all keys in the cache
func (a *App) GetCacheKeys() ([]string, error) {
	if a.cacheDB == nil {
		return nil, fmt.Errorf("cache not initialized")
	}

	var keys []string
	err := a.cacheDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("hurl_cache"))
		if b == nil {
			return fmt.Errorf("cache bucket does not exist")
		}
		return b.ForEach(func(k, v []byte) error {
			keys = append(keys, string(k))
			return nil
		})
	})

	return keys, err
}
