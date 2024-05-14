package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"project01/src/auth"
	"project01/src/db"
	"project01/src/models"
	"project01/src/repositories"
	"project01/src/response"
	"strconv"

	"github.com/gorilla/mux"
)

// NewUser creates a new user
func NewUser(w http.ResponseWriter, r *http.Request) {
	responseBody, err := io.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	err = json.Unmarshal(responseBody, &user)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = user.Prepare("create")
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

	userRepo := repositories.UserRepository{DB: db}

	userID, err := userRepo.Create(&user)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusCreated, userID)
}

// AllUsers returns all users
func AllUsers(w http.ResponseWriter, r *http.Request) {
	db, err := db.New()
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userRepo := repositories.UserRepository{DB: db}

	users, err := userRepo.FindAll()
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, users)
}

// FindUser returns a user
func FindUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	parsedUserID, err := strconv.ParseUint(userID, 10, 64)
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

	userRepo := repositories.UserRepository{DB: db}

	user, err := userRepo.FindByID(parsedUserID)
	if err != nil {
		if err == repositories.ErrNotFound {
			response.ERROR(w, http.StatusNotFound, err)
			return
		}

		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, user)
}

// UpdateUser updates a user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	responseBody, err := io.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	err = json.Unmarshal(responseBody, &user)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

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

	if userIDFromToken != parsedUserID {
		response.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	db, err := db.New()
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userRepo := repositories.UserRepository{DB: db}

	user.ID = parsedUserID

	err = userRepo.Update(&user)
	if err != nil {
		if err == repositories.ErrNotFound {
			response.ERROR(w, http.StatusNotFound, err)
			return
		}

		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// DeleteUser deletes a user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	if userIDFromToken != parsedUserID {
		response.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	db, err := db.New()
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userRepo := repositories.UserRepository{DB: db}

	err = userRepo.Delete(parsedUserID)
	if err != nil {
		if err == repositories.ErrNotFound {
			response.ERROR(w, http.StatusNotFound, err)
			return
		}

		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// FindByFilters returns a list of users filtered by filters
func FindByFilters(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	term := params.Get("term")

	if term == "" {
		response.JSON(w, http.StatusOK, []int{})
		return
	}

	db, err := db.New()
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userRepo := repositories.UserRepository{DB: db}

	users, err := userRepo.FindByFilters(term)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if len(users) == 0 {
		response.JSON(w, http.StatusOK, []int{})
		return
	}

	response.JSON(w, http.StatusOK, users)
}

// FollowUser follows a user
func FollowUser(w http.ResponseWriter, r *http.Request) {
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

	if userIDFromToken == parsedUserID {
		response.ERROR(w, http.StatusBadRequest, errors.New("you can't follow yourself"))
		return
	}

	db, err := db.New()
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userRepo := repositories.UserRepository{DB: db}
	user, err := userRepo.FindByID(parsedUserID)
	if err != nil {
		if err == repositories.ErrNotFound {
			response.ERROR(w, http.StatusNotFound, err)
			return
		}

		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = userRepo.Follow(userIDFromToken, user.ID)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// UnfollowUser unfollows a user
func UnfollowUser(w http.ResponseWriter, r *http.Request) {
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

	if userIDFromToken == parsedUserID {
		response.ERROR(w, http.StatusBadRequest, errors.New("you can't unfollow yourself"))
		return
	}

	db, err := db.New()
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userRepo := repositories.UserRepository{DB: db}
	user, err := userRepo.FindByID(parsedUserID)
	if err != nil {
		if err == repositories.ErrNotFound {
			response.ERROR(w, http.StatusNotFound, err)
			return
		}

		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = userRepo.Unfollow(userIDFromToken, user.ID)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// UserFollowers returns a list of followers from a user
func UserFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	parsedUserID, err := strconv.ParseUint(userID, 10, 64)
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

	userRepo := repositories.UserRepository{DB: db}

	users, err := userRepo.Followers(parsedUserID)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if len(users) == 0 {
		response.JSON(w, http.StatusOK, []int{})
		return
	}

	response.JSON(w, http.StatusOK, users)
}

// UserFollowing returns a list of users that a user is following
func UserFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	parsedUserID, err := strconv.ParseUint(userID, 10, 64)
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

	userRepo := repositories.UserRepository{DB: db}

	users, err := userRepo.Following(parsedUserID)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if len(users) == 0 {
		response.JSON(w, http.StatusOK, []int{})
		return
	}

	response.JSON(w, http.StatusOK, users)
}
