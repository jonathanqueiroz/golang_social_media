package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"project01/src/db"
	"project01/src/models"
	"project01/src/repositories"
	"project01/src/response"
	"strconv"

	"github.com/gorilla/mux"
)

func NewPost(w http.ResponseWriter, r *http.Request) {
	createPost(w, r)
}

func FindPost(w http.ResponseWriter, r *http.Request) {

}

func UpdatePost(w http.ResponseWriter, r *http.Request) {

}

func DeletePost(w http.ResponseWriter, r *http.Request) {

}

func UserPosts(w http.ResponseWriter, r *http.Request) {
}

func LikePost(w http.ResponseWriter, r *http.Request) {

}

func UnlikePost(w http.ResponseWriter, r *http.Request) {

}

func LikesPost(w http.ResponseWriter, r *http.Request) {

}

func NewComment(w http.ResponseWriter, r *http.Request) {
	createPost(w, r)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {

}

func createPost(w http.ResponseWriter, r *http.Request) {
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

	err = post.Prepare()
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.New()
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	postRepo := repositories.PostRepository{DB: db}

	createdPost, err := postRepo.Create(&post)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusCreated, createdPost)
}
