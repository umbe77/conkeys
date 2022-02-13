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

func GetUserStorage(invariantName string) storage.UserStorage {
	var usrProvider storage.UserStorage
	switch invariantName {
	case "memory":
		usrProvider = memory.UserMemoryStorage{}
	case "postgres":
		usrProvider = postgres.PostgresUserStorage{}
	}
	if usrProvider != nil {
		usrProvider.Init()
		return usrProvider
	}
	return nil
}

func GetSecurityStorage(invariantName string) storage.SecurityStorage {
	var securityProvider storage.SecurityStorage
	switch invariantName {
	case "memory":
		securityProvider = memory.SecurityMemoryStorage{}
	}
	if securityProvider != nil {
		securityProvider.Init()
		return securityProvider
	}
	return nil
}
