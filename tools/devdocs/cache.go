package main

import (
	"errors"
	"fmt"
	"os"
	"path"
)

type Cache interface {
	Retrieve(key string) (data []byte, ok bool, err error)
	Store(key string, data []byte) error
}

type MemoryCache struct {
	entries map[string][]byte
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		entries: make(map[string][]byte),
	}
}

// Retrieve implements [Cache].
func (m *MemoryCache) Retrieve(key string) (data []byte, ok bool, err error) {
	data, ok = m.entries[key]
	return data, ok, nil
}

// Store implements [Cache].
func (m *MemoryCache) Store(key string, data []byte) error {
	m.entries[key] = data
	return nil
}

type FileSystemCache struct {
	dir string
}

// Retrieve implements [Cache].
func (f *FileSystemCache) Retrieve(filepath string) (data []byte, ok bool, err error) {
	name := path.Join(f.dir, filepath)
	data, err = os.ReadFile(name)
	if err == nil {
		ok = true
	} else if errors.Is(err, os.ErrNotExist) {
		err = nil
	}

	return data, ok, err
}

// Store implements [Cache].
func (f *FileSystemCache) Store(filepath string, data []byte) error {
	name := path.Join(f.dir, filepath)
	dir := path.Dir(name)

	err := os.MkdirAll(dir, os.ModeDir)
	if err != nil {
		return err
	}

	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func NewFileSystemCache(dir string) *FileSystemCache {
	return &FileSystemCache{
		dir: dir,
	}
}

func userCacheDir() (string, error) {
	d, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	return path.Join(d, "devdocs"), nil
}

var DefaultCache Cache

func init() {
	u, err := userCacheDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: cache broken: %w\n", err)
		DefaultCache = NewMemoryCache()
	} else {
		DefaultCache = NewFileSystemCache(u)
	}
}
