package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"project01/src/auth"
	"project01/src/db"
	"project01/src/models"
	"project01/src/repositories"
	"project01/src/response"

	"golang.org/x/crypto/bcrypt"
)

// Login authenticates a user
func Login(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(requestBody, &user); err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	db, err := db.New()
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userRepo := repositories.UserRepository{DB: db}
	userSaved, err := userRepo.FindByEmail(user.Email)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	passwordErr := bcrypt.CompareHashAndPassword([]byte(userSaved.Password), []byte(user.Password))
	if passwordErr != nil {
		response.ERROR(w, http.StatusUnauthorized, nil)
		return
	}

	token, err := auth.CreateToken(userSaved.ID)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, token)
}
