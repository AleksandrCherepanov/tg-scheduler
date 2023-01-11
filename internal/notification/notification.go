package notification

import (
	"encoding/json"
	"log"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type Storage interface {
	Get() []Notification
	GetByUserId(userId int64) (Notification, error)
	CreateByUserId(userId int64, value string, day time.Weekday, hour int) (Notification, error)
	DeleteByUserId(userId int64) error
	UpdateByUserId(userId int64, value string, day time.Weekday, hour int) (Notification, error)
}

type Notification struct {
	UserId int64        `json:"user_id" validate:"required,min=1"`
	Value  string       `json:"value" validate:"required"`
	Day    time.Weekday `json:"day" validate:"required,max=6"`
	Hour   int          `json:"hour" validate:"required,min=0,max=23"`
}

func NewNotification(userId int64, value string, day time.Weekday, hour int) *Notification {
	return &Notification{
		UserId: userId,
		Value:  value,
		Day:    day,
		Hour:   hour,
	}
}

type NotificationRedisStorage struct {
	sync.Mutex
	client *redis.Client
}

var notificationRedisStorage *NotificationRedisStorage

func GetNotificationRedisStorage(client *redis.Client) Storage {
	if notificationRedisStorage == nil {
		notificationRedisStorage = &NotificationRedisStorage{client: client}
	}

	return notificationRedisStorage
}

func (ns *NotificationRedisStorage) Get() []Notification {
	ns.Lock()

	var cursor uint64
	var keys []string

	for {
		var err error
		keys, cursor, err = ns.client.Scan(cursor, "notification:*", 10).Result()
		if err != nil {
			log.Printf("Get notificaion list error: %s\n", err.Error())
			return []Notification{}
		}

		if cursor == 0 {
			break
		}
	}

	sort.Strings(keys)

	ns.Unlock()

	notifications := make([]Notification, 0, len(keys))
	for _, key := range keys {
		notification, err := ns.getByPrimaryKey(key)
		if err != nil {
			log.Printf("Get notification list error: %s\n", err.Error())
			return []Notification{}
		}

		notifications = append(notifications, notification)
	}

	return notifications
}

func (ns *NotificationRedisStorage) GetByUserId(userId int64) (Notification, error) {
	defer ns.Unlock()
	ns.Lock()

	primaryKey := ns.getPrimaryKey(userId)
	return ns.getByPrimaryKey(primaryKey)
}

func (ns *NotificationRedisStorage) getPrimaryKey(userId int64) string {
	return "notification:" + strconv.FormatInt(userId, 10)
}

func (ns *NotificationRedisStorage) CreateByUserId(
	userId int64,
	value string,
	day time.Weekday,
	hour int,
) (Notification, error) {
	defer ns.Unlock()
	ns.Lock()

	notification := NewNotification(userId, value, day, hour)
	err := ns.save(notification)
	if err != nil {
		log.Printf("Create notification with notification id %d error: %s\n", userId, err.Error())
		return Notification{}, err
	}

	return *notification, nil
}

func (ns *NotificationRedisStorage) UpdateByUserId(
	userId int64,
	value string,
	day time.Weekday,
	hour int,
) (Notification, error) {
	notification, err := ns.GetByUserId(userId)
	if err != nil {
		log.Printf("Update notification with user id %d error: %s\n", userId, err.Error())
		return Notification{}, err
	}

	ns.Lock()
	defer ns.Unlock()

	notification.Value = value
	notification.Day = day
	notification.Hour = hour

	ns.save(&notification)
	if err != nil {
		log.Printf("Update notification with user id %d error: %s\n", userId, err.Error())
		return Notification{}, err
	}

	return notification, nil
}

func (ns *NotificationRedisStorage) DeleteByUserId(userId int64) error {
	ns.Lock()
	defer ns.Unlock()

	primaryKey := ns.getPrimaryKey(userId)
	err := ns.client.Del(primaryKey).Err()
	if err != nil {
		log.Printf("Delete notification with %d error: %s\n", userId, err.Error())
		return err
	}

	return err
}

func (ns *NotificationRedisStorage) save(notification *Notification) error {
	jsonNotification, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Save notification with user id %d error: %s\n", notification.UserId, err.Error())
		return err
	}

	primaryKey := ns.getPrimaryKey(notification.UserId)

	err = ns.client.Set(primaryKey, jsonNotification, 0).Err()
	if err != nil {
		log.Printf("Save notification with user id %d error: %s\n", notification.UserId, err.Error())
		return err
	}

	return nil
}

func (ns *NotificationRedisStorage) getByPrimaryKey(key string) (Notification, error) {
	notification := Notification{}
	jsonNotification, err := ns.client.Get(key).Bytes()
	if err != nil {
		log.Printf("Get notification with key %s error: %s\n", key, err.Error())
		return notification, err
	}

	err = json.Unmarshal(jsonNotification, &notification)
	if err != nil {
		log.Printf("Get notification with key %s error: %s\n", key, err.Error())
		return notification, err
	}

	return notification, nil
}
