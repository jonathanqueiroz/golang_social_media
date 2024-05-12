package repositories

import (
	"database/sql"
	"errors"
	"project01/src/models"
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
	query := `SELECT id, name, email, birthdate, created_at FROM users`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Birthdate, &user.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) FindByID(id uint64) (*models.User, error) {
	query := `SELECT id, name, email, birthdate, created_at FROM users WHERE id = $1`
	rows := r.DB.QueryRow(query, id)

	var user models.User
	if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Birthdate, &user.CreatedAt); err != nil {
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

func (r *UserRepository) FindByName(name string) (*models.User, error) {
	query := `SELECT id, name, email, birthdate, created_at FROM users WHERE name LIKE $1`
	name = "%" + name + "%"
	rows := r.DB.QueryRow(query, name)

	var user models.User
	if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Birthdate, &user.CreatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}
