package scheduler

import (
	"fmt"
	"time"

	"github.com/AleksandrCherepanov/tg-scheduler/internal/notification"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/user"
)

const seconds_in_hour = 3600
const tgNotification_day = time.Thursday
const tgNotification_utc_hour = 12

var notificator *Notificator
var sender *asyncSender

type Notificator struct {
	userStorage *user.UserStorage
	sender      *asyncSender
}

type tgNotification struct {
	userId  int64
	message string
}

func GetNotificator() *Notificator {
	if notificator == nil {
		notificator = &Notificator{}
		notificator.userStorage = user.GetUserStorage()
		notificator.sender = getSender()
	}

	return notificator
}

type asyncSender struct {
	input chan tgNotification
}

func getSender() *asyncSender {
	if sender == nil {
		sender = &asyncSender{}
		sender.input = make(chan tgNotification)
	}

	return sender
}

func (n *Notificator) StartNotification() {
	currentUTCTime := time.Now().UTC().Unix()
	secondsPassedHour := currentUTCTime % seconds_in_hour
	secondsTillNewHour := seconds_in_hour - secondsPassedHour
	time.Sleep(time.Duration(secondsTillNewHour) * time.Second)

	for range time.Tick(time.Hour) {
		n.notify()
	}
}

func (n *Notificator) notify() {
	if n.isCorrectDay() && n.isCorrectTime() {
		users := n.userStorage.GetAllUsers()
		for _, user := range users {
			text, err := notification.GetNotificationStorage().Get(user.Id)
			tn := tgNotification{
				userId: user.Id,
			}
			if err != nil {
				tn.message = err.Error()
			} else {
				tn.message = text.(string)
			}

			go n.sender.send()
			n.sender.input <- tn
		}
	}
}

func (s *asyncSender) send() {
	for {
		select {
		case n := <-s.input:
			fmt.Println(n.userId)
			fmt.Println(n.message)
		}
	}
}

func (n *Notificator) isCorrectDay() bool {
	return time.Now().UTC().Weekday() == tgNotification_day
}

func (n *Notificator) isCorrectTime() bool {
	return time.Now().UTC().Hour() >= tgNotification_utc_hour
}
