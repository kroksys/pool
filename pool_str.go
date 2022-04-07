package pool

import (
	"sync"
)

// Creates new pool with any type provided
func NewPoolStr[T any]() *poolStr[T] {
	return &poolStr[T]{
		lock:    &sync.RWMutex{},
		storage: make(map[string]T),
	}
}

type poolStr[T any] struct {

	// holds the data
	storage map[string]T

	// thread safe interactions with map
	lock *sync.RWMutex
}

// Adds object to pool and returns index ID that can be used
// later to retreive the object
func (p *poolStr[T]) Put(id string, object T) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.storage[id] = object
}

// Gets object from storage. If value does not exists it will
// return empty struct/object/interface
func (p *poolStr[T]) Get(id string) T {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.storage[id]
}

// The same as Get with difference that it return boolean
// weather the value exists in storage.
func (p *poolStr[T]) GetOk(id string) (T, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	_, ok := p.storage[id]
	return p.storage[id], ok
}

// Deletes element from storage. If value does not exist it does nothin
func (p *poolStr[T]) Delete(id string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	delete(p.storage, id)
}

// Executes a function on stored element
func (p *poolStr[T]) Exec(id string, fn func(T)) {
	fn(p.Get(id))
}

// Executes a function on stored element and returns modified element
func (p *poolStr[T]) Map(id string, fn func(T) T) T {
	return fn(p.Get(id))
}
