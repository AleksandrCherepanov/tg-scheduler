package user

import (
	"testing"
)

func TestCreateUser(t *testing.T) {
	userStorage := GetUserStorage()
	userStorage.CreateUser(1, "test1")
	userStorage.CreateUser(2, "test2")

	result := userStorage.GetAllUsers()
	checkUser(t, result, 1, "test1")
	checkUser(t, result, 2, "test2")
}

func checkUser(t *testing.T, result map[int64]User, id int64, name string) {
	if _, ok := result[id]; !ok {
		t.Fatalf("user %d not found", id)
	}

	user := result[id]
	if user.Id != id {
		t.Fatalf("invalid. expected user id %d actual user id %d", id, user.Id)
	}
	if user.Name != name {
		t.Fatalf("invalid. expected user name %s actual user name %s", name, user.Name)
	}
}
