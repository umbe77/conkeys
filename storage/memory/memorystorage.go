package memory

import (
	"conkeys/storage"
	"errors"
	"fmt"
	"strings"
)

var c = make(map[string]storage.Value)

type MemoryStorage struct {

}

func (m MemoryStorage) Init() {}

func (m MemoryStorage) Get(path string) (storage.Value, error) {
	value, ok := c[path]
	if ok {
		return value, nil
	}
	return storage.Value{},  errors.New(fmt.Sprintf("%s key not present in storage", path))
}

func (m MemoryStorage) GetKeys(pathSearch string) (map[string]storage.Value, error) {
	result := make(map[string]storage.Value)
	found := false
	for k := range c {
		if strings.HasPrefix(k, pathSearch) {
			found = true
			result[k] = c[k]
		}
	}
	if found {
		return result, nil
	}
	return result, errors.New(fmt.Sprintf("No keys found for %s", pathSearch))
}

func (m MemoryStorage) GetAllKeys() map[string]storage.Value {
	return c
}

func (m MemoryStorage) Put(path string, value storage.Value) {
	c[path] = value
}