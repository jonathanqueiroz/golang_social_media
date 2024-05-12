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

func DeleteUser(w http.ResponseWriter, r *http.Request) {
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

func FindUserByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["name"]

	db, err := db.New()
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userRepo := repositories.UserRepository{DB: db}

	user, err := userRepo.FindByName(userName)
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
