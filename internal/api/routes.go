package api

import (
	"github.com/AleksandrCherepanov/tg-scheduler/internal/notification"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/storage"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/user"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	userAPI := GetUserAPI(user.GetUserStorage(storage.GetRedisClient()))
	r.HandleFunc("/api/user", userAPI.GetUserList).Methods("GET")
	r.HandleFunc("/api/user/{id:\\d+}", userAPI.GetUser).Methods("GET")
	r.HandleFunc("/api/user", userAPI.CreateUser).Methods("POST")
	r.HandleFunc("/api/user/{id:\\d+}", userAPI.UpdateUser).Methods("PUT")
	r.HandleFunc("/api/user/{id:\\d+}", userAPI.DeleteUser).Methods("DELETE")

	notificationAPI := GetNotificationAPI(
		notification.GetNotificationRedisStorage(
			storage.GetRedisClient(),
		),
	)
	r.HandleFunc("/api/notification", notificationAPI.GetNotificationList).Methods("GET")
	r.HandleFunc("/api/notification/{user_id:\\d+}", notificationAPI.GetNotification).Methods("GET")
	r.HandleFunc("/api/notification", notificationAPI.CreateNotification).Methods("POST")
	r.HandleFunc("/api/notification/{user_id:\\d+}", notificationAPI.UpdateNotification).Methods("PUT")
	r.HandleFunc("/api/notification/{user_id:\\d+}", notificationAPI.DeleteNotification).Methods("DELETE")
}
