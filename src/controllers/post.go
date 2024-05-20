package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"project01/src/auth"
	"project01/src/models"
	"project01/src/repositories"
	"project01/src/response"
	"project01/src/websocket"
	"strconv"

	"github.com/gorilla/mux"
)

type PostController struct {
	PostRepo         repositories.PostRepositoryInterface
	NotificationRepo repositories.NotificationRepositoryInterface
}

func NewPostController(db *sql.DB) *PostController {
	return &PostController{
		PostRepo:         repositories.NewPostRepository(db),
		NotificationRepo: repositories.NewNotificationRepository(db),
	}
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

	createdPost, err := pc.PostRepo.Create(&post)
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

	posts, err := pc.PostRepo.PostsFollowedUsers(userIDFromToken)
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

	post, err := pc.PostRepo.FindByID(parsedPostID, userIDFromToken)
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

	post, err := pc.PostRepo.FindByID(parsedPostID, userIDFromToken)
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
	err = pc.PostRepo.Update(post)
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

	post, err := pc.PostRepo.FindByID(parsedPostID, userIDFromToken)
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

	err = pc.PostRepo.Delete(post.ID)
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

	posts, err := pc.PostRepo.FindByAuthorID(parsedUserID, userIDFromToken)
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

	post, err := pc.PostRepo.FindByID(parsedPostID, userIDFromToken)
	if err != nil {
		if err == repositories.ErrNotFound {
			response.ERROR(w, http.StatusNotFound, err)
			return
		}

		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	newLikeInserted, err := pc.PostRepo.LikePost(post.ID, userIDFromToken)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if newLikeInserted && post.AuthorID != userIDFromToken {
		notification := models.Notification{
			UserID:       post.AuthorID,
			Type:         "like",
			SourceUserID: userIDFromToken,
			SourcePostID: &post.ID,
		}

		err = pc.NotificationRepo.CreateOrUpdate(notification)
		if err != nil {
			log.Println(err)
		}

		websocket.SendNotification(post.AuthorID, notification)
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

	post, err := pc.PostRepo.FindByID(parsedPostID, userIDFromToken)
	if err != nil {
		if err == repositories.ErrNotFound {
			response.ERROR(w, http.StatusNotFound, err)
			return
		}

		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = pc.PostRepo.UnlikePost(post.ID, userIDFromToken)
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

	likes, err := pc.PostRepo.LikesPost(parsedPostID)
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
