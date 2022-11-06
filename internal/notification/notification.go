package notification

import (
	"fmt"
	"sync"
)

type NotificationStorage struct {
	sync.Mutex
	notifications map[int64]string
}

var notificationStorage *NotificationStorage

func GetNotificationStorage() *NotificationStorage {
	if notificationStorage == nil {
		notificationStorage = &NotificationStorage{}
		notificationStorage.notifications = make(map[int64]string)
	}

	return notificationStorage
}

func (ns *NotificationStorage) Set(userId int64, notificaiton string) {
	ns.Lock()
	defer ns.Unlock()

	ns.notifications[userId] = notificaiton
}

func (ns *NotificationStorage) Get(userId int64) (interface{}, error) {
	ns.Lock()
	defer ns.Unlock()

	if notification, ok := ns.notifications[userId]; ok {
		return notification, nil
	}

	return nil, fmt.Errorf("notification for user=%d not found", userId)
}
