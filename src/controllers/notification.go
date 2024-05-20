package controllers

import (
	"database/sql"
	"net/http"
	"project01/src/auth"
	"project01/src/models"
	"project01/src/repositories"
	"project01/src/response"
	"strconv"

	"github.com/gorilla/mux"
)

type NotificationController struct {
	NotificationRepo repositories.NotificationRepositoryInterface
}

func NewNotificationController(db *sql.DB) *NotificationController {
	return &NotificationController{
		NotificationRepo: repositories.NewNotificationRepository(db),
	}
}

// FindAllNotifications returns all notifications for a user
func (nc *NotificationController) FindAllNotifications(w http.ResponseWriter, r *http.Request) {
	userIDFromToken, err := auth.ExtractUserID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	notifications, err := nc.NotificationRepo.FindAll(userIDFromToken)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if len(notifications) == 0 {
		response.JSON(w, http.StatusOK, []models.Notification{})
		return
	}

	response.JSON(w, http.StatusOK, notifications)
}

// FindNotificationByID returns a notification by ID
func (nc *NotificationController) FindNotificationByID(w http.ResponseWriter, r *http.Request) {
	userIDFromToken, err := auth.ExtractUserID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	notification, err := nc.NotificationRepo.FindByID(userIDFromToken, parsedID)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, notification)
}

// DeleteNotification deletes a notification
func (nc *NotificationController) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	userIDFromToken, err := auth.ExtractUserID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	notification, err := nc.NotificationRepo.FindByID(userIDFromToken, parsedID)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if userIDFromToken != notification.UserID {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	err = nc.NotificationRepo.Delete(userIDFromToken, parsedID)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// MarkNotificationAsRead marks a notification as read
func (nc *NotificationController) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	userIDFromToken, err := auth.ExtractUserID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	notification, err := nc.NotificationRepo.FindByID(userIDFromToken, parsedID)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if userIDFromToken != notification.UserID {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	err = nc.NotificationRepo.MarkIDAsRead(userIDFromToken, parsedID)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// MarkAllNotificationsAsRead marks all notifications as read
func (nc *NotificationController) MarkAllNotificationsAsRead(w http.ResponseWriter, r *http.Request) {
	userIDFromToken, err := auth.ExtractUserID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	err = nc.NotificationRepo.MarkAllAsRead(userIDFromToken)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
