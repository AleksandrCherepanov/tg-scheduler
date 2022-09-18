package scheduler

import (
	"fmt"
	"time"

	"github.com/AleksandrCherepanov/tg-scheduler/internal/user"
)

func Notify() {
	userStorage := user.GetUserStorage()

	for range time.Tick(time.Second * 5) {
		users := userStorage.GetAllUsers()
		for _, user := range users {
			fmt.Println(user.Id)
			fmt.Println(user.Name)
		}
	}
}
