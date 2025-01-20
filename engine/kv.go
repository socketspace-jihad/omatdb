package engine

import (
	"errors"
	"sync"
)

const (
	KeyNotExistsErr string = "key is not exists"
)

type kv struct {
	mtx  *sync.RWMutex
	data map[string]any
}

func NewKVStore() *kv {
	return &kv{
		mtx:  &sync.RWMutex{},
		data: make(map[string]any),
	}
}

func (s *kv) Store(key string, val any) error {
	_, ok := s.data[key]
	if ok {
		return errors.New("key is exists, use StoreOW instead for overwriting")
	}
	s.mtx.Lock()
	s.data[key] = val
	s.mtx.Unlock()
	return nil
}

func (s *kv) StoreOW(key string, val any) error {
	s.mtx.Lock()
	s.data[key] = val
	s.mtx.Unlock()
	return nil
}

func (s *kv) keyExists(key string) bool {
	_, ok := s.data[key]
	return ok
}

func (s *kv) Get(key string) (any, error) {
	if !s.keyExists(key) {
		return nil, errors.New(KeyNotExistsErr)
	}
	return s.data[key], nil
}

func (s *kv) Update(key string, val any) error {
	if !s.keyExists(key) {
		return errors.New(KeyNotExistsErr)
	}
	s.mtx.Lock()
	s.data[key] = val
	s.mtx.Unlock()
	return nil
}

func (s *kv) Delete(key string) error {
	if !s.keyExists(key) {
		return errors.New(KeyNotExistsErr)
	}
	s.mtx.Lock()
	delete(s.data, key)
	s.mtx.Unlock()
	return nil
}
