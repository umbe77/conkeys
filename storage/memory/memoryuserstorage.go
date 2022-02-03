package memory

import (
	"conkeys/storage"
	"fmt"
	"strings"
)

var users = make(map[string]storage.User)

type UserMemoryStorage struct{}

func (u UserMemoryStorage) Init() {}

func (u UserMemoryStorage) Get(userName string) (storage.User, error) {
	user, ok := users[userName]
	if ok {
		return user, nil
	}
	return storage.User{}, fmt.Errorf("%s user not present in storage", userName)
}

func (u UserMemoryStorage) GetUsers(query string) ([]storage.User, error) {
	result := make([]storage.User, 1, 1)
	found := false
	for k, v := range users {
		if strings.HasPrefix(k, query) {
			found = true
			result = append(result, v)
		}
	}
	if found {
		return result, nil
	}
	return result, fmt.Errorf("no users found for %s", query)

}

func (u UserMemoryStorage) Add (usr storage.User) {
	users[usr.UserName] = usr
}

func (u UserMemoryStorage) Update (usr storage.User) {
	users[usr.UserName] = usr
}

func (u UserMemoryStorage) Delete (userName string) {
	delete(users, userName)
}

func (u UserMemoryStorage) SetPassword(userName string, password string) {
	usr, ok := users[userName]
	if ok {
		usr.Password = password
	}
}
