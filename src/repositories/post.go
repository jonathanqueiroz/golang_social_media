package repositories

import (
	"database/sql"
	"project01/src/models"
)

type PostRepositoryInterface interface {
	Create(post *models.Post) (*models.Post, error)
	FindByAuthorID(authorID uint64, currentUserID uint64) ([]models.Post, error)
	FindByID(id uint64, currentUserID uint64) (*models.Post, error)
	Update(post *models.Post) error
	Delete(id uint64) error
	LikePost(postID, userID uint64) error
	UnlikePost(postID, userID uint64) error
	LikesPost(postID uint64) ([]models.User, error)
	PostsFollowedUsers(userID uint64) ([]models.Post, error)
}

func NewPostRepository(db *sql.DB) PostRepositoryInterface {
	return &PostRepository{DB: db}
}

type PostRepository struct {
	DB *sql.DB
}

// Create creates a new post in the database.
func (r *PostRepository) Create(post *models.Post) (*models.Post, error) {
	query := `INSERT INTO posts (content, author_id, parent_id) VALUES ($1, $2, $3) RETURNING id`

	var id uint64

	err := r.DB.QueryRow(query, post.Content, post.AuthorID, post.ParentID).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return r.FindByID(id, post.AuthorID)
}

// FindByAuthorID retrieves posts by the author's ID from the database.
func (r *PostRepository) FindByAuthorID(authorID uint64, currentUserID uint64) ([]models.Post, error) {
	var posts []models.Post

	query := `SELECT posts.*, users.name AS author_name, users.username, COUNT(likes.id) AS total_likes,
		(SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = posts.id AND user_id = $2)) AS current_user_liked,
		(SELECT COUNT(*) FROM posts AS replies WHERE replies.parent_id = posts.id) AS total_replies
		FROM posts
		LEFT JOIN users ON users.id = posts.author_id
		LEFT JOIN likes ON likes.post_id = posts.id
		WHERE posts.author_id = $1 AND posts.parent_id IS NULL
		GROUP BY posts.id, users.name, users.username
		ORDER BY posts.created_at DESC`
	rows, err := r.DB.Query(query, authorID, currentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(
			&post.ID,
			&post.ParentID,
			&post.AuthorID,
			&post.Content,
			&post.CreatedAt,
			&post.AuthorName,
			&post.Username,
			&post.TotalLikes,
			&post.CurrentUserLiked,
			&post.TotalReplies,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// FindByID retrieves a post by its ID from the database.
func (r *PostRepository) FindByID(id uint64, currentUserID uint64) (*models.Post, error) {
	var post models.Post

	query := `SELECT posts.*, users.name AS author_name, users.username, COUNT(likes.id) AS total_likes,
		(SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = posts.id AND user_id = $2)) AS current_user_liked,
		(SELECT COUNT(*) FROM posts AS replies WHERE replies.parent_id = posts.id) AS total_replies
		FROM posts
		LEFT JOIN users ON users.id = posts.author_id
		LEFT JOIN likes ON likes.post_id = posts.id
		WHERE posts.id = $1
		GROUP BY posts.id, users.name, users.username`
	err := r.DB.QueryRow(query, id, currentUserID).Scan(
		&post.ID,
		&post.ParentID,
		&post.AuthorID,
		&post.Content,
		&post.CreatedAt,
		&post.AuthorName,
		&post.Username,
		&post.TotalLikes,
		&post.CurrentUserLiked,
		&post.TotalReplies,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	post.Replies, err = r.findRepliesByParentID(post.ID, currentUserID)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

// Update updates the content of a post in the database.
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

// Delete deletes a post from the database.
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

// LikePost adds a like to a post in the database.
func (r *PostRepository) LikePost(postID, userID uint64) error {
	query := `INSERT INTO likes (post_id, user_id) VALUES ($1, $2) ON CONFLICT (post_id, user_id) DO NOTHING`
	_, err := r.DB.Exec(query, postID, userID)
	if err != nil {
		return err
	}

	return nil
}

// UnlikePost removes a like from a post in the database.
func (r *PostRepository) UnlikePost(postID, userID uint64) error {
	query := `DELETE FROM likes WHERE post_id = $1 AND user_id = $2`
	_, err := r.DB.Exec(query, postID, userID)
	if err != nil {
		return err
	}

	return nil
}

// LikesPost retrieves the users who liked a post from the database.
func (r *PostRepository) LikesPost(postID uint64) ([]models.User, error) {
	var users []models.User

	query := `SELECT users.id, name, username, avatar_url, bio, users.created_at FROM likes LEFT JOIN users ON users.id = likes.user_id WHERE likes.post_id = $1 ORDER BY likes.id`
	rows, err := r.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.AvatarURL, &user.Bio, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// PostsFollowedUsers retrieves posts from the users that the current user is following.
func (r *PostRepository) PostsFollowedUsers(userID uint64) ([]models.Post, error) {
	var posts []models.Post

	query := `SELECT posts.*, users.name AS author_name, users.username, COUNT(likes.id) AS total_likes,
		(SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = posts.id AND user_id = $1)) AS current_user_liked,
		(SELECT COUNT(*) FROM posts AS replies WHERE replies.parent_id = posts.id) AS total_replies
		FROM posts
		LEFT JOIN users ON users.id = posts.author_id
		LEFT JOIN likes ON likes.post_id = posts.id
		WHERE posts.author_id IN (SELECT user_id FROM followers WHERE follower_id = $1) AND posts.parent_id IS NULL
		GROUP BY posts.id, users.name, users.username
		ORDER BY posts.created_at DESC`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(
			&post.ID,
			&post.ParentID,
			&post.AuthorID,
			&post.Content,
			&post.CreatedAt,
			&post.AuthorName,
			&post.Username,
			&post.TotalLikes,
			&post.CurrentUserLiked,
			&post.TotalReplies,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// findRepliesByParentID retrieves the replies to a post from the database.
func (r *PostRepository) findRepliesByParentID(parentID uint64, currentUserID uint64) ([]models.Post, error) {
	var posts []models.Post

	query := `SELECT posts.*, users.name AS author_name, users.username, COUNT(likes.id) AS total_likes,
		(SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = posts.id AND user_id = $2)) AS current_user_liked,
		(SELECT COUNT(*) FROM posts AS replies WHERE replies.parent_id = posts.id) AS total_replies
		FROM posts
		LEFT JOIN users ON users.id = posts.author_id
		LEFT JOIN likes ON likes.post_id = posts.id
		WHERE posts.parent_id = $1
		GROUP BY posts.id, users.name, users.username
		ORDER BY total_likes DESC, posts.created_at ASC`
	rows, err := r.DB.Query(query, parentID, currentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(
			&post.ID,
			&post.ParentID,
			&post.AuthorID,
			&post.Content,
			&post.CreatedAt,
			&post.AuthorName,
			&post.Username,
			&post.TotalLikes,
			&post.CurrentUserLiked,
			&post.TotalReplies,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
