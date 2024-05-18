package controllers

import (
	"net/http"
	"project01/src/auth"
	"project01/src/models"
	"project01/src/repositories"
	"project01/src/response"
)

type ProfileController struct {
	UserRepo *repositories.UserRepository
	PostRepo *repositories.PostRepository
}

func NewProfileController(userRepo *repositories.UserRepository, postRepo *repositories.PostRepository) *ProfileController {
	return &ProfileController{
		UserRepo: userRepo,
		PostRepo: postRepo,
	}
}

func (pc *ProfileController) GetProfile(w http.ResponseWriter, r *http.Request) {
	userIDFromToken, err := auth.ExtractUserID(r)
	if err != nil {
		response.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	type Result struct {
		Profile *models.Profile
		Err     error
	}

	resultChan := make(chan Result, 4)

	go func() {
		user, err := pc.UserRepo.FindByID(userIDFromToken)
		resultChan <- Result{Profile: &models.Profile{User: *user}, Err: err}
	}()

	go func() {
		posts, err := pc.PostRepo.FindByAuthorID(userIDFromToken, userIDFromToken)
		resultChan <- Result{Profile: &models.Profile{Posts: posts}, Err: err}
	}()

	go func() {
		followers, err := pc.UserRepo.Followers(userIDFromToken)
		resultChan <- Result{Profile: &models.Profile{FollowersCount: len(followers)}, Err: err}
	}()

	go func() {
		following, err := pc.UserRepo.Following(userIDFromToken)
		resultChan <- Result{Profile: &models.Profile{FollowingCount: len(following)}, Err: err}
	}()

	profile := &models.Profile{}

	for i := 0; i < cap(resultChan); i++ {
		var userBlank models.User

		result := <-resultChan
		if result.Err != nil {
			response.ERROR(w, http.StatusInternalServerError, result.Err)
			return
		}
		if result.Profile.User != userBlank {
			profile.User = result.Profile.User
		}
		if result.Profile.Posts != nil {
			profile.Posts = result.Profile.Posts
		}
		if result.Profile.FollowersCount != 0 {
			profile.FollowersCount = result.Profile.FollowersCount
		}
		if result.Profile.FollowingCount != 0 {
			profile.FollowingCount = result.Profile.FollowingCount
		}
	}

	response.JSON(w, http.StatusOK, profile)
}
