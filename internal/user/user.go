package user

import (
	"encoding/json"
	"log"
	"sort"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
)

type Storage interface {
	Get() []User
	GetById(id int64) (User, error)
	Create(id int64, name string) (User, error)
	Delete(id int64) error
	Update(id int64, name string) (User, error)
}

type User struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type UserRedisStorage struct {
	sync.Mutex
	client *redis.Client
}

var userStorage *UserRedisStorage

func GetUserStorage(client *redis.Client) Storage {
	if userStorage == nil {
		userStorage = &UserRedisStorage{client: client}
	}

	return userStorage
}

func (us *UserRedisStorage) Get() []User {
	us.Lock()

	var cursor uint64
	var keys []string

	for {
		var err error
		keys, cursor, err = us.client.Scan(cursor, "user:*", 10).Result()
		if err != nil {
			log.Printf("Get user list error: %s\n", err.Error())
			return []User{}
		}

		if cursor == 0 {
			break
		}
	}

	sort.Strings(keys)

	us.Unlock()

	users := make([]User, 0, len(keys))
	for _, key := range keys {
		user, err := us.getByPrimaryKey(key)
		if err != nil {
			log.Printf("Get user list error: %s\n", err.Error())
			return []User{}
		}

		users = append(users, user)
	}

	return users
}

func (us *UserRedisStorage) GetById(id int64) (User, error) {
	defer us.Unlock()
	us.Lock()

	primaryKey := us.getPrimaryKey(id)
	return us.getByPrimaryKey(primaryKey)
}

func (us *UserRedisStorage) getByPrimaryKey(key string) (User, error) {
	user := User{}
	jsonUser, err := us.client.Get(key).Bytes()
	if err != nil {
		log.Printf("Get user with key %s error: %s\n", key, err.Error())
		return user, err
	}

	err = json.Unmarshal(jsonUser, &user)
	if err != nil {
		log.Printf("Get user with key %s error: %s\n", key, err.Error())
		return user, err
	}

	return user, nil
}

func (us *UserRedisStorage) Create(id int64, name string) (User, error) {
	defer us.Unlock()
	us.Lock()

	user := User{Id: id, Name: name}
	err := us.save(user)
	if err != nil {
		log.Printf("Create user with id %d error: %s\n", id, err.Error())
		return User{}, err
	}

	return user, nil
}

func (us *UserRedisStorage) Update(id int64, name string) (User, error) {
	user, err := us.GetById(id)
	if err != nil {
		log.Printf("Update user with id %d error: %s\n", id, err.Error())
		return User{}, err
	}

	us.Lock()
	defer us.Unlock()

	user.Name = name
	us.save(user)
	if err != nil {
		log.Printf("Update user with id %d error: %s\n", id, err.Error())
		return User{}, err
	}

	return user, nil
}

func (us *UserRedisStorage) Delete(id int64) error {
	us.Lock()
	defer us.Unlock()

	primaryKey := us.getPrimaryKey(id)
	err := us.client.Del(primaryKey).Err()
	if err != nil {
		log.Printf("Delete user with %d error: %s\n", id, err.Error())
		return err
	}

	return err
}

func (us *UserRedisStorage) save(user User) error {
	jsonUser, err := json.Marshal(user)
	if err != nil {
		log.Printf("Save user with %d error: %s\n", user.Id, err.Error())
		return err
	}

	primaryKey := us.getPrimaryKey(user.Id)

	err = us.client.Set(primaryKey, jsonUser, 0).Err()
	if err != nil {
		log.Printf("Save user with %d error: %s\n", user.Id, err.Error())
		return err
	}

	return nil
}

func (us *UserRedisStorage) getPrimaryKey(id int64) string {
	return "user:" + strconv.FormatInt(id, 10)
}
