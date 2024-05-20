package routes

import (
	"database/sql"
	"net/http"
	"project01/src/controllers"
)

func notificationRoutes(db *sql.DB) []Route {
	notificationController := controllers.NewNotificationController(db)

	return []Route{
		{
			URI:          "/notifications",
			Method:       http.MethodGet,
			Function:     notificationController.FindAllNotifications,
			AuthRequired: true,
		},
		{
			URI:          "/notifications/{id}",
			Method:       http.MethodGet,
			Function:     notificationController.FindNotificationByID,
			AuthRequired: true,
		},
		{
			URI:          "/notifications",
			Method:       http.MethodDelete,
			Function:     notificationController.DeleteNotification,
			AuthRequired: true,
		},
		{
			URI:          "/notifications/{id}/read",
			Method:       http.MethodPut,
			Function:     notificationController.MarkNotificationAsRead,
			AuthRequired: true,
		},
		{
			URI:          "/notifications/read-all",
			Method:       http.MethodPut,
			Function:     notificationController.MarkAllNotificationsAsRead,
			AuthRequired: true,
		},
	}
}
