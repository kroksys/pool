package pool

import (
	"sync"
)

// Creates new pool with any type provided
func NewPool[T any]() *pool[T] {
	return &pool[T]{
		lock:    &sync.RWMutex{},
		storage: make(map[uint64]T),
	}
}

type pool[T any] struct {

	// index of stored element is incremented on each Put request
	// Limitations:
	// id uint64 maximum number is 18446744073709551615 so this is a limit how many
	// elements can pool store
	index uint64

	// holds the data
	storage map[uint64]T

	// thread safe interactions with map
	lock *sync.RWMutex
}

// Adds object to pool and returns index ID that can be used
// later to retreive the object
func (p *pool[T]) Put(object T) uint64 {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.storage[p.index] = object
	return p.index
}

// Gets object from storage. If value does not exists it will
// return empty struct/object/interface
func (p *pool[T]) Get(id uint64) T {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.storage[id]
}

// The same as Get with difference that it return boolean
// weather the value exists in storage.
func (p *pool[T]) GetOk(id uint64) (T, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	_, ok := p.storage[id]
	return p.storage[id], ok
}

// Deletes element from storage. If value does not exist it does nothin
func (p *pool[T]) Delete(id uint64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	delete(p.storage, id)
}

// Executes a function on stored element
func (p *pool[T]) Exec(id uint64, fn func(T)) {
	fn(p.Get(id))
}

// Executes a function on stored element and returns modified element
func (p *pool[T]) Map(id uint64, fn func(T) T) T {
	return fn(p.Get(id))
}
