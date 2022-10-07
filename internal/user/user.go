package user

import (
	"fmt"
	"sync"
)

type User struct {
	Id   int64
	Name string
}

type UserStorage struct {
	sync.Mutex
	userList map[int64]User
}

var userStorage *UserStorage

func GetUserStorage() *UserStorage {
	if userStorage == nil {
		userStorage = &UserStorage{}
		userStorage.userList = make(map[int64]User)
	}

	return userStorage
}

func (userStorage *UserStorage) GetAllUsers() map[int64]User {
	return userStorage.userList
}

func (userStorage *UserStorage) CreateUser(id int64, name string) User {
	userStorage.Lock()
	defer userStorage.Unlock()

	user := User{}
	user.Id = id
	user.Name = name
	userStorage.userList[id] = user

	return user
}

func (userStorage *UserStorage) GetUserById(id int64) (User, error) {
	userStorage.Lock()
	defer userStorage.Unlock()

	user, ok := userStorage.userList[id]
	if !ok {
		return User{}, fmt.Errorf("user with id=%d not found", id)
	}

	return user, nil
}
