package scheduler

import (
	"fmt"
	"time"

	"github.com/AleksandrCherepanov/go_telegram/pkg/telegram/client"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/notification"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/user"
)

const secondsInHour = 3600
const tickPeriod = time.Hour
const notificationDay = time.Friday
const notificationUtcHour = 16

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

func sleepBeforeNewHour() {
	currentUTCTime := time.Now().UTC().Unix()
	secondsPassedHour := currentUTCTime % secondsInHour
	secondsTillNewHour := secondsInHour - secondsPassedHour
	time.Sleep(time.Duration(secondsTillNewHour) * time.Second)
}

func (n *Notificator) StartNotification() {
	sleepBeforeNewHour()
	n.notify()
	for range time.Tick(tickPeriod) {
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
			r := client.NewTelegramResponse(n.userId, n.message, false)
			res, err := r.Send()
			fmt.Println(err)
			fmt.Println(string(res.([]byte)))
		}
	}
}

func (n *Notificator) isCorrectDay() bool {
	return time.Now().UTC().Weekday() == notificationDay
}

func (n *Notificator) isCorrectTime() bool {
	currentHour := time.Now().UTC().Hour()
	return currentHour >= notificationUtcHour && currentHour < notificationUtcHour+1
}
