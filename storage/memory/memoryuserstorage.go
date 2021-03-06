package memory

import (
	"conkeys/storage"
	"errors"
	"fmt"
	"strings"
)

var users = make(map[string]storage.User)
var passwords = make(map[string]string)

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

func (u UserMemoryStorage) Add(usr storage.User) error {
	users[usr.UserName] = usr
	passwords[usr.UserName] = ""
	return nil
}

func (u UserMemoryStorage) Update(usr storage.User) error {
	users[usr.UserName] = usr
	return nil
}

func (u UserMemoryStorage) Delete(userName string) error {
	delete(users, userName)
	return nil
}

func (u UserMemoryStorage) SetPassword(userName string, password string) error {
	_, ok := passwords[userName]
	if ok {
		passwords[userName] = password
	}
	return nil
}

func (u UserMemoryStorage) GetPassword(userName string) (string, storage.User, error) {
	pwd, okPwd := passwords[userName]
	user, okUsr := users[userName]
	if !okPwd || !okUsr {
		return "", storage.User{}, errors.New("User not present")
	}
	return pwd, user, nil
}
