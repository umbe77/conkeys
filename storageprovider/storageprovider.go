package storageprovider

import (
	"conkeys/storage"
	"conkeys/storage/memory"
	"conkeys/storage/mongodb"
	"conkeys/storage/postgres"
)

func GetKeyStorage(invariantName string) storage.KeyStorage {
	var stgProvider storage.KeyStorage
	switch invariantName {
	case "memory":
		stgProvider = memory.MemoryStorage{}
	case "mongodb":
		stgProvider = mongodb.MongoStorage{}
	case "postgres":
		stgProvider = postgres.PostgresStorage{}
	}
	if stgProvider != nil {
		stgProvider.Init()
		return stgProvider
	}
	return nil
}
