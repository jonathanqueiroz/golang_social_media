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

type UserController struct {
	Repo repositories.UserRepositoryInterface
}

func NewUserController(db *sql.DB) *UserController {
	return &UserController{
		Repo: repositories.NewUserRepository(db),
	}
}

func (uc *UserController) GetRepo() *repositories.UserRepository {
	return uc.Repo.(*repositories.UserRepository)
}

// NewUser creates a new user
func (uc *UserController) NewUser(w http.ResponseWriter, r *http.Request) {
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

	userID, err := uc.Repo.Create(&user)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusCreated, userID)
}

// AllUsers returns all users
func (uc *UserController) AllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uc.Repo.FindAll()
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, users)
}

// FindUser returns a user
func (uc *UserController) FindUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	parsedUserID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user, err := uc.Repo.FindByID(parsedUserID)
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
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	user.ID = parsedUserID

	err = uc.Repo.Update(&user)
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
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	err = uc.Repo.Delete(parsedUserID)
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
func (uc *UserController) FindByFilters(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	term := params.Get("term")

	if term == "" {
		response.JSON(w, http.StatusOK, []int{})
		return
	}

	users, err := uc.Repo.FindByFilters(term)
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
func (uc *UserController) FollowUser(w http.ResponseWriter, r *http.Request) {
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

	user, err := uc.Repo.FindByID(parsedUserID)
	if err != nil {
		if err == repositories.ErrNotFound {
			response.ERROR(w, http.StatusNotFound, err)
			return
		}

		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = uc.Repo.Follow(userIDFromToken, user.ID)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// UnfollowUser unfollows a user
func (uc *UserController) UnfollowUser(w http.ResponseWriter, r *http.Request) {
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

	user, err := uc.Repo.FindByID(parsedUserID)
	if err != nil {
		if err == repositories.ErrNotFound {
			response.ERROR(w, http.StatusNotFound, err)
			return
		}

		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = uc.Repo.Unfollow(userIDFromToken, user.ID)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// UserFollowers returns a list of followers from a user
func (uc *UserController) UserFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	parsedUserID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	users, err := uc.Repo.Followers(parsedUserID)
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
func (uc *UserController) UserFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	parsedUserID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	users, err := uc.Repo.Following(parsedUserID)
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
