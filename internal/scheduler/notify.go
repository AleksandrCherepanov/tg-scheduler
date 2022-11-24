package scheduler

import (
	"fmt"
	"time"

	"github.com/AleksandrCherepanov/go_telegram/pkg/telegram/client"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/notification"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/storage"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/user"
)

const secondsInHour = 3600
const tickPeriod = time.Hour
const NotificationDay = time.Monday
const NotificationUtcHour = 4

var notificator *Notificator
var sender *asyncSender

type Notificator struct {
	userStorage         user.Storage
	notificationStorage notification.Storage
	sender              *asyncSender
}

type tgNotification struct {
	userId  int64
	message string
}

func GetNotificator() *Notificator {
	if notificator == nil {
		notificator = &Notificator{}
		notificator.userStorage = user.GetUserStorage(storage.GetRedisClient())
		notificator.notificationStorage = notification.GetNotificationRedisStorage(
			storage.GetRedisClient(),
		)
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
	users := n.userStorage.Get()
	for _, user := range users {
		notification, err := n.notificationStorage.GetByUserId(user.Id)
		tn := tgNotification{
			userId: user.Id,
		}
		if err != nil {
			tn.message = err.Error()
		} else {
			tn.message = notification.Value
		}

		if n.isCorrectDay(notification) && n.isCorrectTime(notification) {
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

func (n *Notificator) isCorrectDay(ntf notification.Notification) bool {
	return time.Now().UTC().Weekday() == ntf.Day
}

func (n *Notificator) isCorrectTime(ntf notification.Notification) bool {
	currentHour := time.Now().UTC().Hour()
	return currentHour >= ntf.Hour && currentHour < ntf.Hour+1
}
