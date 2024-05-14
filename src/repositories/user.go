package repositories

import (
	"database/sql"
	"errors"
	"project01/src/models"
	"strings"
)

type UserRepository struct {
	DB *sql.DB
}

var ErrNotFound = errors.New("not found")

func (r *UserRepository) Create(user *models.User) (*models.User, error) {
	query := `INSERT INTO users (name, email, password, birthdate) VALUES ($1, $2, $3, $4) RETURNING id`

	var id uint64

	err := r.DB.QueryRow(query, user.Name, user.Email, user.Password, user.Birthdate).Scan(&id)

	if err != nil {
		return nil, err
	}

	return r.FindByID(id)
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	query := `SELECT id, name, email, username, avatar_url, bio, birthdate, created_at FROM users ORDER BY name ASC`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Username,
			&user.AvatarURL,
			&user.Bio,
			&user.Birthdate,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) FindByID(id uint64) (*models.User, error) {
	query := `SELECT id, name, email, username, avatar_url, bio, birthdate, created_at FROM users WHERE id = $1`
	rows := r.DB.QueryRow(query, id)

	var user models.User
	if err := rows.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Username,
		&user.AvatarURL,
		&user.Bio,
		&user.Birthdate,
		&user.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := `UPDATE users SET name = $1, email = $2, birthdate = $3 WHERE id = $4`
	result, err := r.DB.Exec(query, user.Name, user.Email, user.Birthdate, user.ID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *UserRepository) Delete(id uint64) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.DB.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *UserRepository) FindByFilters(term string) ([]models.User, error) {
	term = strings.ToLower(term)

	query := `SELECT id, name, email, username, avatar_url, bio, birthdate, created_at
				FROM users
        WHERE LOWER(name) LIKE $1 OR LOWER(email) LIKE $1 OR LOWER(username) LIKE $1
        ORDER BY name ASC`
	rows, err := r.DB.Query(query, "%"+term+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Username,
			&user.AvatarURL,
			&user.Bio,
			&user.Birthdate,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `SELECT id, password FROM users WHERE email = $1`
	rows := r.DB.QueryRow(query, email)

	var user models.User
	if err := rows.Scan(&user.ID, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Follow(followerID, userID uint64) error {
	query := `INSERT INTO followers (follower_id, user_id) VALUES ($1, $2) ON CONFLICT (follower_id, user_id) DO NOTHING`
	_, err := r.DB.Exec(query, followerID, userID)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Unfollow(followerID, userID uint64) error {
	query := `DELETE FROM followers WHERE follower_id = $1 AND user_id = $2`
	_, err := r.DB.Exec(query, followerID, userID)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Followers(userID uint64) ([]models.User, error) {
	query := `SELECT users.id, users.name, users.username, users.avatar_url, users.bio
		FROM followers
		LEFT JOIN users ON users.id = followers.follower_id
		WHERE followers.user_id = $1 ORDER BY followers.created_at DESC`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.AvatarURL, &user.Bio); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) Following(userID uint64) ([]models.User, error) {
	query := `SELECT users.id, users.name, users.username, users.avatar_url, users.bio
		FROM followers
		LEFT JOIN users ON users.id = followers.user_id
		WHERE followers.follower_id = $1 ORDER BY followers.created_at DESC`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.AvatarURL, &user.Bio); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
