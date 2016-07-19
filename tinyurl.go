package tinyurl

import (
	"bytes"
	"math/rand"
	"sync"
)

const (
	charset = "0123456789abcdefghjilmnopqrsutvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	charsetb = []byte(charset)
	charsetl = len(charsetb)
)

// HashFunc define funtion to hash url
type HashFunc func(value []byte, hashlen int) string

func naiveHash(value []byte, hashlen int) string {
	var b bytes.Buffer

	for b.Len() < hashlen {
		b.WriteByte(charsetb[rand.Intn(charsetl)])
	}

	return b.String()
}

// NewTinyURL return default tiny url struct
// Using in-memory backend, app built-in naiveHash and check saved url
func NewTinyURL() *TinyURL {
	mb := NewMemBackend()
	return NewCustomTinyURL(mb, naiveHash, 8, true)
}

// NewCustomTinyURL accepts custom params
func NewCustomTinyURL(backend Backend, hashFn HashFunc, hashLen int, checkSaved bool) *TinyURL {
	tu := &TinyURL{}
	tu.backend = backend
	tu.hashFn = hashFn
	tu.checkdu = checkSaved
	tu.hashLen = hashLen
	return tu
}

// TinyURL struct
type TinyURL struct {
	mu      sync.Mutex
	backend Backend
	hashFn  HashFunc
	hashLen int
	checkdu bool
}

// Shorten given URL, assume given urls are valid...
func (tu *TinyURL) Shorten(url string) (string, error) {

	if tu.checkdu {
		if v, ok := tu.backend.CheckURL(url); ok {
			return v, nil
		}
	}

	// Trick: use counter value to hash instead of origin url
	c := append([]byte{}, byte(tu.backend.GetCount()))
	v := ""
	for {
		h := tu.hashFn(c, tu.hashLen)
		err := tu.backend.SetURL(h, url)
		if err != nil {
			continue
		}
		v = h
		break
	}

	return v, nil
}

// Recover given hash to origin url
func (tu *TinyURL) Recover(hash string) (string, error) {
	if v, ok := tu.backend.GetURL(hash); ok {
		return v, nil
	}
	return "", NotExistError{Hash: hash}
}

// SetHashLen update hashlen
func (tu *TinyURL) SetHashLen(l int) {
	tu.mu.Lock()
	defer tu.mu.Unlock()
	tu.hashLen = l
}
