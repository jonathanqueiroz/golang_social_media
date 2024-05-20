package repositories

import (
	"database/sql"
	"project01/src/models"
)

type NotificationRepositoryInterface interface {
	FindAll(userID uint64) ([]models.Notification, error)
	FindByID(userID, id uint64) (models.Notification, error)
	CreateOrUpdate(notification models.Notification) error
	Delete(userID, id uint64) error
	MarkIDAsRead(userID, id uint64) error
	MarkAllAsRead(userID uint64) error
}

type NotificationRepository struct {
	DB *sql.DB
}

func NewNotificationRepository(db *sql.DB) NotificationRepositoryInterface {
	return &NotificationRepository{DB: db}
}

func (r *NotificationRepository) FindAll(userID uint64) ([]models.Notification, error) {
	var notifications []models.Notification

	query := `SELECT DISTINCT notifications.*, users.name, users.username, users.avatar_url, posts.content AS post_content,
		CASE
			WHEN type = 'new_follower' THEN GREATEST((SELECT COUNT(*) FROM followers WHERE user_id = notifications.user_id) - 1, 0)
			WHEN type = 'like' THEN GREATEST((SELECT COUNT(*) FROM likes WHERE post_id = notifications.source_post_id) - 1, 0)
			ELSE 0 END AS others_total
		FROM notifications
		LEFT JOIN users ON notifications.source_user_id = users.id
		LEFT JOIN posts ON notifications.source_post_id = posts.id
		WHERE notifications.user_id = $1
		GROUP BY notifications.id, users.name, users.username, users.avatar_url, posts.content
		ORDER BY created_at DESC`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var notification models.Notification
		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Type,
			&notification.SourceUserID,
			&notification.SourcePostID,
			&notification.IsRead,
			&notification.CreatedAt,
			&notification.UpdatedAt,
			&notification.Name,
			&notification.Username,
			&notification.AvatarURL,
			&notification.PostContent,
			&notification.OthersTotal,
		)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (r *NotificationRepository) FindByID(userID, id uint64) (models.Notification, error) {
	var notification models.Notification

	query := `SELECT notifications.*, users.name, users.username, users.avatar_url, posts.content AS post_content
		FROM notifications
		LEFT JOIN users ON notifications.source_user_id = users.id
		LEFT JOIN posts ON notifications.source_post_id = posts.id
		WHERE user_id = $1 AND notifications.id = $2`
	row := r.DB.QueryRow(query, userID, id)

	err := row.Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Type,
		&notification.SourceUserID,
		&notification.SourcePostID,
		&notification.IsRead,
		&notification.CreatedAt,
		&notification.UpdatedAt,
		&notification.Name,
		&notification.Username,
		&notification.AvatarURL,
		&notification.PostContent,
	)
	if err != nil {
		return models.Notification{}, err
	}

	return notification, nil
}

func (r *NotificationRepository) CreateOrUpdate(notification models.Notification) error {
	result, err := r.DB.Exec(`
			UPDATE notifications
			SET source_user_id = $1, updated_at = CURRENT_TIMESTAMP, is_read = FALSE
			WHERE user_id = $3 AND type = $2
		`, notification.SourceUserID, notification.Type, notification.UserID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		_, err := r.DB.Exec(`
					INSERT INTO notifications (user_id, type, source_user_id, source_post_id)
					VALUES ($1, $2, $3, $4)
			`, notification.UserID, notification.Type, notification.SourceUserID, notification.SourcePostID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *NotificationRepository) Delete(userID, id uint64) error {
	_, err := r.DB.Exec(`DELETE FROM notifications WHERE user_id = $1 AND id = $2`, userID, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) MarkIDAsRead(userID, id uint64) error {
	_, err := r.DB.Exec(`UPDATE notifications SET is_read = TRUE WHERE user_id = $1 AND id = $2`, userID, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) MarkAllAsRead(userID uint64) error {
	_, err := r.DB.Exec(`UPDATE notifications SET is_read = TRUE WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	return nil
}
