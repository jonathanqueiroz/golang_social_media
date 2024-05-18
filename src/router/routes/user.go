package routes

import (
	"database/sql"
	"net/http"
	"project01/src/controllers"
)

func userRoutes(db *sql.DB) []Route {
	userController := controllers.NewUserController(db)

	return []Route{
		{
			URI:          "/users",
			Method:       http.MethodPost,
			Function:     userController.NewUser,
			AuthRequired: false,
		},
		{
			URI:          "/users",
			Method:       http.MethodGet,
			Function:     userController.FindByFilters,
			AuthRequired: true,
		},
		{
			URI:          "/users/{id}",
			Method:       http.MethodGet,
			Function:     userController.FindUser,
			AuthRequired: true,
		},
		{
			URI:          "/users/{id}",
			Method:       http.MethodPut,
			Function:     userController.UpdateUser,
			AuthRequired: true,
		},
		{
			URI:          "/users/{id}",
			Method:       http.MethodDelete,
			Function:     userController.DeleteUser,
			AuthRequired: true,
		},
		{
			URI:          "/users/{id}/follow",
			Method:       http.MethodPost,
			Function:     userController.FollowUser,
			AuthRequired: true,
		},
		{
			URI:          "/users/{id}/unfollow",
			Method:       http.MethodPost,
			Function:     userController.UnfollowUser,
			AuthRequired: true,
		},
		{
			URI:          "/users/{id}/followers",
			Method:       http.MethodGet,
			Function:     userController.UserFollowers,
			AuthRequired: true,
		},
		{
			URI:          "/users/{id}/following",
			Method:       http.MethodGet,
			Function:     userController.UserFollowing,
			AuthRequired: true,
		},
	}
}
