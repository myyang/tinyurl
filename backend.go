package tinyurl

import (
	"fmt"
	"sync"
)

// Backend define backend to store encoded and mapping relation
// for further usage.
type Backend interface {
	GetCount() int64
	GetURL(hash string) (url string, ok bool)
	SetURL(hash, url string) error
	CheckURL(url string) (hash string, ok bool)
}

// CollisionError indicates hash collision
type CollisionError struct {
	Hash, New, Old string
}

func (e CollisionError) Error() string {
	return fmt.Sprintf("Hash(%v) collision, New:% v Saved:% v", e.Hash, e.New, e.Old)
}

// NotExistError indicates no such hash value
type NotExistError struct {
	Hash string
}

func (e NotExistError) Error() string {
	return fmt.Sprintf("Hash(%v) not exists")
}

// NewMemBackend return new membackend struct
func NewMemBackend() *MemBackend {
	mb := &MemBackend{}
	mb.urls = make(map[string]string)
	mb.revurls = make(map[string]string)
	return mb
}

// MemBackend is memory stored backend
type MemBackend struct {
	mu      sync.RWMutex
	urls    map[string]string
	revurls map[string]string
	counter int64
}

// GetCount return counter number for hashing
func (mb *MemBackend) GetCount() int64 {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	mb.counter++
	return mb.counter
}

// GetURL check if there is any stored hash value
func (mb *MemBackend) GetURL(hash string) (url string, ok bool) {
	url, ok = mb.urls[hash]
	return url, ok
}

// SetURL save new hash mapping
func (mb *MemBackend) SetURL(hash, url string) error {
	if v, ok := mb.urls[hash]; ok && v != url {
		return CollisionError{Hash: hash, New: url, Old: v}
	}
	mb.mu.Lock()
	defer mb.mu.Unlock()
	mb.urls[hash] = url
	mb.revurls[url] = hash
	return nil
}

// CheckURL return hash if url exists
func (mb *MemBackend) CheckURL(url string) (hash string, ok bool) {
	hash, ok = mb.revurls[url]
	return hash, ok
}
