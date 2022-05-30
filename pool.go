package pool

import (
	"sync"
)

// Creates new pool with any type provided
func NewPool[T comparable]() *Pool[T] {
	return &Pool[T]{
		index:   1,
		lock:    &sync.RWMutex{},
		storage: make(map[uint64]T),
	}
}

type Pool[T comparable] struct {

	// index of stored element is incremented on each Put request
	// Limitations:
	// id uint64 maximum number is 18446744073709551615 so this is a limit how many
	// elements can pool store
	// 0 is reserved for none - not found
	index uint64

	// holds the data
	storage map[uint64]T

	// thread safe interactions with map
	lock *sync.RWMutex
}

// Adds object to pool and returns index ID that can be used
// later to retreive the object
func (p *Pool[T]) Put(object T) uint64 {
	p.lock.Lock()
	defer p.lock.Unlock()
	i := p.index
	p.index++
	p.storage[i] = object
	return i
}

// Gets object from storage. If value does not exists it will
// return empty struct/object/interface
func (p *Pool[T]) Get(id uint64) T {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.storage[id]
}

// The same as Get with difference that it return boolean
// weather the value exists in storage.
func (p *Pool[T]) GetOk(id uint64) (T, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	_, ok := p.storage[id]
	return p.storage[id], ok
}

// Deletes element from storage. If value does not exist it does nothin
func (p *Pool[T]) Delete(id uint64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	delete(p.storage, id)
}

// Executes a function on stored element
func (p *Pool[T]) Exec(id uint64, fn func(T)) {
	fn(p.Get(id))
}

// Executes a function on stored element and returns modified element
func (p *Pool[T]) Map(id uint64, fn func(T) T) T {
	return fn(p.Get(id))
}

func (p *Pool[T]) Find(obj T) (uint64, *T) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	for k, v := range p.storage {
		if v == obj {
			return k, &v
		}
	}
	return 0, nil
}

// Executes a function for each element in pool
func (p *Pool[T]) Each(fn func(T)) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	for _, v := range p.storage {
		fn(v)
	}
}

// Read lock for manual work with data
func (p *Pool[T]) Lock() {
	p.lock.Lock()
}

// Read unlock for manual work with data
func (p *Pool[T]) Unlock() {
	p.lock.Unlock()
}

// Returns storage for manual work with data.
func (p *Pool[T]) Data() map[uint64]T {
	return p.storage
}
