package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"project01/src/auth"
	"project01/src/models"
	"project01/src/repositories"
	"project01/src/response"
	"strconv"

	"github.com/gorilla/mux"
)

type PostController struct {
	Repo repositories.PostRepositoryInterface
}

func NewPostController(db *sql.DB) *PostController {
	return &PostController{
		Repo: repositories.NewPostRepository(db),
	}
}

func (pc *PostController) GetRepo() *repositories.PostRepository {
	return pc.Repo.(*repositories.PostRepository)
}

// NewPost creates a new post
func (pc *PostController) NewPost(w http.ResponseWriter, r *http.Request) {
	responseBody, err := io.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var post models.Post
	err = json.Unmarshal(responseBody, &post)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	vars := mux.Vars(r)
	postID := vars["id"]

	if postID != "" {
		parsedpostID, err := strconv.ParseUint(postID, 10, 64)

		if err != nil {
			response.ERROR(w, http.StatusBadRequest, err)
			return
		}

		post.ParentID = &parsedpostID
	}

	userIDFromToken, err := auth.ExtractUserID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	post.AuthorID = userIDFromToken

	err = post.Prepare()
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	createdPost, err := pc.Repo.Create(&post)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusCreated, createdPost)
}

// PostsFollowedUsers returns a list of posts from followed users
func (pc *PostController) PostsFollowedUsers(w http.ResponseWriter, r *http.Request) {
	userIDFromToken, err := auth.ExtractUserID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	posts, err := pc.Repo.PostsFollowedUsers(userIDFromToken)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if len(posts) == 0 {
		response.JSON(w, http.StatusOK, []int{})
		return
	}

	response.JSON(w, http.StatusOK, posts)
}

// FindPost returns a post
func (pc *PostController) FindPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	parsedPostID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userIDFromToken, err := auth.ExtractUserID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	post, err := pc.Repo.FindByID(parsedPostID, userIDFromToken)
	if err != nil {
		if err == repositories.ErrNotFound {
			response.ERROR(w, http.StatusNotFound, err)
			return
		}

		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, post)
}

// UpdatePost updates a post
func (pc *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	parsedPostID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responseBody, err := io.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var updatedPost models.Post
	err = json.Unmarshal(responseBody, &updatedPost)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userIDFromToken, err := auth.ExtractUserID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	post, err := pc.Repo.FindByID(parsedPostID, userIDFromToken)
	if err != nil {
		if err == repositories.ErrNotFound {
			response.ERROR(w, http.StatusNotFound, err)
			return
		}

		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if userIDFromToken != post.AuthorID {
		response.ERROR(w, http.StatusForbidden, errors.New("unauthorized"))
		return
	}

	post.Content = updatedPost.Content

	post.Prepare()
	err = pc.Repo.Update(post)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// DeletePost deletes a post
func (pc *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	parsedPostID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userIDFromToken, err := auth.ExtractUserID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	post, err := pc.Repo.FindByID(parsedPostID, userIDFromToken)
	if err != nil {
		if err == repositories.ErrNotFound {
			response.ERROR(w, http.StatusNotFound, err)
			return
		}

		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if userIDFromToken != post.AuthorID {
		response.ERROR(w, http.StatusForbidden, errors.New("unauthorized"))
		return
	}

	err = pc.Repo.Delete(post.ID)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// UserPosts returns a list of posts from a user
func (pc *PostController) UserPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	parsedUserID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userIDFromToken, err := auth.ExtractUserID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	posts, err := pc.Repo.FindByAuthorID(parsedUserID, userIDFromToken)
	if err != nil {
		if err == repositories.ErrNotFound {
			response.ERROR(w, http.StatusNotFound, err)
			return
		}

		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if len(posts) == 0 {
		response.JSON(w, http.StatusOK, []int{})
		return
	}

	response.JSON(w, http.StatusOK, posts)
}

// LikePost likes a post
func (pc *PostController) LikePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	parsedPostID, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userIDFromToken, err := auth.ExtractUserID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	post, err := pc.Repo.FindByID(parsedPostID, userIDFromToken)
	if err != nil {
		if err == repositories.ErrNotFound {
			response.ERROR(w, http.StatusNotFound, err)
			return
		}

		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = pc.Repo.LikePost(post.ID, userIDFromToken)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// UnlikePost removes a like from a post
func (pc *PostController) UnlikePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	parsedPostID, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userIDFromToken, err := auth.ExtractUserID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	post, err := pc.Repo.FindByID(parsedPostID, userIDFromToken)
	if err != nil {
		if err == repositories.ErrNotFound {
			response.ERROR(w, http.StatusNotFound, err)
			return
		}

		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = pc.Repo.UnlikePost(post.ID, userIDFromToken)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// LikesPost returns a list of users who liked a post
func (pc *PostController) LikesPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	parsedPostID, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	likes, err := pc.Repo.LikesPost(parsedPostID)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if len(likes) == 0 {
		response.JSON(w, http.StatusOK, []int{})
		return
	}

	response.JSON(w, http.StatusOK, likes)
}
