package repositories

import (
	"database/sql"
	"project01/src/models"
)

type PostRepository struct {
	DB *sql.DB
}

func (r *PostRepository) Create(post *models.Post) (*models.Post, error) {
	query := `INSERT INTO posts (content, author_id, parent_id) VALUES ($1, $2, $3) RETURNING id`

	var id uint64

	err := r.DB.QueryRow(query, post.Content, post.AuthorID, post.ParentID).Scan(&id)

	if err != nil {
		return nil, err
	}

	return r.FindByID(id)
}

func (r *PostRepository) FindByAuthor(authorID uint64) ([]models.Post, error) {
	query := `SELECT id, content, created_at, author_id FROM posts WHERE author_id = $1`
	rows, err := r.DB.Query(query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Content, &post.CreatedAt, &post.AuthorID); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) FindByID(id uint64) (*models.Post, error) {
	query := `SELECT id, content, created_at, author_id FROM posts WHERE id = $1`
	rows := r.DB.QueryRow(query, id)

	var post models.Post
	if err := rows.Scan(&post.ID, &post.Content, &post.CreatedAt, &post.AuthorID); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) Update(post *models.Post) error {
	query := `UPDATE posts SET content = $1 WHERE id = $2`
	result, err := r.DB.Exec(query, post.Content, post.ID)

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

func (r *PostRepository) Delete(id uint64) error {
	query := `DELETE FROM posts WHERE id = $1`
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
