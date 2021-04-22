package storageprovider

import (
	"conkeys/storage"
	"conkeys/storage/memory"
	"conkeys/storage/mongodb"
)

func GetKeyStorage(invariantName string) storage.KeyStorage {
	var stgProvider storage.KeyStorage
	switch invariantName {
	case "memory":
		stgProvider = memory.MemoryStorage{}
	case "mongodb":
		stgProvider = mongodb.MongoStorage{}
	}
	if stgProvider != nil {
		stgProvider.Init()
		return stgProvider
	}
	return nil
}
