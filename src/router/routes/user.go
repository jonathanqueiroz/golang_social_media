package routes

import (
	"net/http"
	"project01/src/controllers"
)

var userRoutes = []Route{
	{
		URI:          "/users",
		Method:       http.MethodPost,
		Function:     controllers.NewUser,
		AuthRequired: false,
	},
	{
		URI:          "/users",
		Method:       http.MethodGet,
		Function:     controllers.FindByFilters,
		AuthRequired: true,
	},
	{
		URI:          "/users/{id}",
		Method:       http.MethodGet,
		Function:     controllers.FindUser,
		AuthRequired: true,
	},
	{
		URI:          "/users/{id}",
		Method:       http.MethodPut,
		Function:     controllers.UpdateUser,
		AuthRequired: true,
	},
	{
		URI:          "/users/{id}",
		Method:       http.MethodDelete,
		Function:     controllers.DeleteUser,
		AuthRequired: true,
	},
	{
		URI:          "/users/{id}/follow",
		Method:       http.MethodPost,
		Function:     controllers.FollowUser,
		AuthRequired: true,
	},
	{
		URI:          "/users/{id}/unfollow",
		Method:       http.MethodPost,
		Function:     controllers.UnfollowUser,
		AuthRequired: true,
	},
	{
		URI:          "/users/{id}/followers",
		Method:       http.MethodGet,
		Function:     controllers.UserFollowers,
		AuthRequired: true,
	},
	{
		URI:          "/users/{id}/following",
		Method:       http.MethodGet,
		Function:     controllers.UserFollowing,
		AuthRequired: true,
	},
}
