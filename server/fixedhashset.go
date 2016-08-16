package server

import (
	"bytes"
	"sync"
)

type fixedHashSet struct {
	shards [64]struct {
		rLock  sync.RWMutex
		wLock  sync.Mutex
		offset int
		hashes [1024][16]byte
	}
}

// Contains returns true if hash is in the hashset
func (s *fixedHashSet) Contains(hash []byte) bool {
	// Find relevant shard
	shard := &s.shards[(hash[0]^hash[1]^hash[2]^hash[3])%64]

	// Aqcuire read lock
	shard.rLock.Lock()
	defer shard.rLock.Unlock()

	// Check if the hash exists
	for i := 0; i < 1024; i++ {
		if bytes.Equal(shard.hashes[i][:], hash) {
			return true
		}
	}
	return false
}

// Insert hash
func (s *fixedHashSet) Insert(hash []byte) {
	// Find relevant shard
	shard := &s.shards[(hash[0]^hash[1]^hash[2]^hash[3])%64]

	// Acquire write lock (upgrading control)
	shard.wLock.Lock()
	defer shard.wLock.Unlock()

	// Insert hash
	copy(shard.hashes[shard.offset][:], hash)
	shard.offset = (shard.offset + 1) % 1024
}
