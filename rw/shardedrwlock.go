// Package shardedrwlock provides a sharded read-write lock mechanism to reduce lock contention
// in concurrent applications by distributing locks across multiple shards based on the hash of a key.
package cxlockrw

/*
#cgo LDFLAGS: -lpthread
#include <pthread.h>
#include <stdlib.h>

// Initializes a pthread read-write lock.
void rwlock_init(pthread_rwlock_t *lock) {
    pthread_rwlock_init(lock, NULL);
}

// Destroys a pthread read-write lock.
void rwlock_destroy(pthread_rwlock_t *lock) {
    pthread_rwlock_destroy(lock);
}

// Acquires a read lock on a pthread read-write lock.
void rwlock_rlock(pthread_rwlock_t *lock) {
    pthread_rwlock_rdlock(lock);
}

// Releases a read lock on a pthread read-write lock.
void rwlock_runlock(pthread_rwlock_t *lock) {
    pthread_rwlock_unlock(lock);
}

// Acquires a write lock on a pthread read-write lock.
void rwlock_lock(pthread_rwlock_t *lock) {
    pthread_rwlock_wrlock(lock);
}

// Releases a write lock on a pthread read-write lock.
void rwlock_unlock(pthread_rwlock_t *lock) {
    pthread_rwlock_unlock(lock);
}
*/
import "C"
import (
	"hash/fnv"
	"runtime"
	"unsafe"
)

// RWLockShard represents a single shard containing a POSIX read-write lock.
type RWLockShard struct {
	rwlock C.pthread_rwlock_t
}

// init initializes the shard's read-write lock.
func (shard *RWLockShard) init() {
	C.rwlock_init(&shard.rwlock)
}

// destroy destroys the shard's read-write lock.
func (shard *RWLockShard) destroy() {
	C.rwlock_destroy(&shard.rwlock)
}

// rlock acquires a read lock for the shard.
func (shard *RWLockShard) rlock() {
	C.rwlock_rlock(&shard.rwlock)
}

// runlock releases a read lock for the shard.
func (shard *RWLockShard) runlock() {
	C.rwlock_runlock(&shard.rwlock)
}

// lock acquires a write lock for the shard.
func (shard *RWLockShard) lock() {
	C.rwlock_lock(&shard.rwlock)
}

// unlock releases a write lock for the shard.
func (shard *RWLockShard) unlock() {
	C.rwlock_unlock(&shard.rwlock)
}

// ShardedRWLock provides a set of sharded read-write locks to reduce lock contention.
type ShardedRWLock struct {
	shards []RWLockShard
}

// NewShardedRWLock creates a new ShardedRWLock with a specified number of shards.
func NewShardedRWLock(numShards int) *ShardedRWLock {
	lock := &ShardedRWLock{
		shards: make([]RWLockShard, numShards),
	}
	for i := range lock.shards {
		lock.shards[i].init()
	}
	return lock
}

// Close cleans up resources used by the ShardedRWLock.
func (lock *ShardedRWLock) Close() {
	for i := range lock.shards {
		lock.shards[i].destroy()
	}
}

// getShard selects the appropriate shard based on the hash of a key.
func (lock *ShardedRWLock) getShard(key string) *RWLockShard {
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(key))
	hash := hasher.Sum32()
	return &lock.shards[hash%uint32(len(lock.shards))]
}

// RLock acquires a read lock for the shard corresponding to the provided key.
func (lock *ShardedRWLock) RLock(key string) {
	shard := lock.getShard(key)
	shard.rlock()
}

// RUnlock releases a read lock for the shard corresponding to the provided key.
func (lock *ShardedRWLock) RUnlock(key string) {
	shard := lock.getShard(key)
	shard.runlock()
}

// Lock acquires a write lock for the shard corresponding to the provided key.
func (lock *ShardedRWLock) Lock(key string) {
	shard := lock.getShard(key)
	shard.lock()
}

// Unlock releases a write lock for the shard corresponding to the provided key.
func (lock *ShardedRWLock) Unlock(key string) {
	shard := lock.getShard(key)
	shard.unlock()
}
