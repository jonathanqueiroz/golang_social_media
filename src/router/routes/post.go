package routes

import (
	"net/http"
	"project01/src/controllers"
)

var postRoutes = []Route{
	{
		URI:          "/posts",
		Method:       http.MethodPost,
		Function:     controllers.NewPost,
		AuthRequired: true,
	},
	{
		URI:          "/posts/{id}",
		Method:       http.MethodGet,
		Function:     controllers.FindPost,
		AuthRequired: true,
	},
	{
		URI:          "/posts/{id}",
		Method:       http.MethodPut,
		Function:     controllers.UpdatePost,
		AuthRequired: true,
	},
	{
		URI:          "/posts/{id}",
		Method:       http.MethodDelete,
		Function:     controllers.DeletePost,
		AuthRequired: true,
	},
	{
		URI:          "/users/{id}/posts",
		Method:       http.MethodGet,
		Function:     controllers.UserPosts,
		AuthRequired: true,
	},
	{
		URI:          "/posts/{id}/like",
		Method:       http.MethodPost,
		Function:     controllers.LikePost,
		AuthRequired: true,
	},
	{
		URI:          "/posts/{id}/unlike",
		Method:       http.MethodPost,
		Function:     controllers.UnlikePost,
		AuthRequired: true,
	},
	{
		URI:          "/posts/{id}/likes",
		Method:       http.MethodGet,
		Function:     controllers.LikesPost,
		AuthRequired: true,
	},
	{
		URI:          "/posts/{id}/comments",
		Method:       http.MethodPost,
		Function:     controllers.NewComment,
		AuthRequired: true,
	},
	{
		URI:          "/posts/{id}/comments",
		Method:       http.MethodDelete,
		Function:     controllers.DeleteComment,
		AuthRequired: true,
	},
}
