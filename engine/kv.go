package engine

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

const (
	KeyNotExistsErr string = "key is not exists"
)

type KVStorage struct {
	mtx     *sync.RWMutex
	data    map[string]any
	dataDir string
}

func NewKVStore(dataDir string) *KVStorage {
	return &KVStorage{
		mtx:     &sync.RWMutex{},
		data:    make(map[string]any),
		dataDir: dataDir,
	}
}

func (s *KVStorage) Store(key string, val any) error {
	_, ok := s.data[key]
	if ok {
		return errors.New("key is exists, use StoreOW instead for overwriting")
	}
	s.mtx.Lock()
	s.data[key] = val
	s.mtx.Unlock()
	return nil
}

func (s *KVStorage) StoreOW(key string, val any) error {
	s.mtx.Lock()
	s.data[key] = val
	s.mtx.Unlock()
	return nil
}

func (s *KVStorage) keyExists(key string) bool {
	_, ok := s.data[key]
	return ok
}

func (s *KVStorage) Get(key string) (any, error) {
	if !s.keyExists(key) {
		return nil, errors.New(KeyNotExistsErr)
	}
	return s.data[key], nil
}

func (s *KVStorage) Update(key string, val any) error {
	if !s.keyExists(key) {
		return errors.New(KeyNotExistsErr)
	}
	s.mtx.Lock()
	s.data[key] = val
	s.mtx.Unlock()
	return nil
}

func (s *KVStorage) Delete(key string) error {
	if !s.keyExists(key) {
		return errors.New(KeyNotExistsErr)
	}
	s.mtx.Lock()
	delete(s.data, key)
	s.mtx.Unlock()
	return nil
}

func (k *KVStorage) Flush() error {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	k.mtx.Lock()
	data, err := json.Marshal(k.data)
	if err != nil {
		k.mtx.Unlock()
		return err
	}
	k.mtx.Unlock()
	writer.Write(data)
	writer.Close()
	return os.WriteFile(fmt.Sprintf("%s/omatdb.gz", k.dataDir), buf.Bytes(), 0666)
}

func (k *KVStorage) Load() error {
	fi, err := os.Open(fmt.Sprintf("%s/omatdb.gz", k.dataDir))
	if err != nil {
		return err
	}
	defer fi.Close()
	reader, err := gzip.NewReader(fi)
	if err != nil {
		return err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &k.data); err != nil {
		return err
	}

	return nil
}
