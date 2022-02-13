package memory

import (
	"conkeys/storage"
	"fmt"
	"strings"
)

var c = make(map[string]storage.Value)
var cEncrypted = make(map[string]string)

type MemoryStorage struct {
}

func (m MemoryStorage) Init() {}

func (m MemoryStorage) Get(path string) (storage.Value, error) {
	value, ok := c[path]
	if ok {
		return value, nil
	}
	return storage.Value{}, fmt.Errorf("%s key not present in storage", path)
}

func (m MemoryStorage) GetEncrypted(path string) (storage.Value, error) {
	value, ok := cEncrypted[path]
	if ok {
		return storage.Value{
			T: storage.Crypted,
			V: value,
		}, nil
	}
	return storage.Value{}, fmt.Errorf("%s key not present in storage", path)
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
	return result, fmt.Errorf("no keys found for %s", pathSearch)
}

func (m MemoryStorage) GetAllKeys() map[string]storage.Value {
	return c
}

func (m MemoryStorage) Put(path string, value storage.Value) {
	c[path] = value
}

func (m MemoryStorage) PutEncrypted(path string, maskedValue storage.Value, encryptedValue string) {
	c[path] = maskedValue
	cEncrypted[path] = encryptedValue
	// c[path] = storage.Value{
	// 	T: value.T,
	// 	V: "********",
	// }
	// cEncrypted[path] = fmt.Sprintf("%s", value.V)
}

func (m MemoryStorage) Delete(path string) {
	delete(c, path)
}
