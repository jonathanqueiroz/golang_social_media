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
		AuthRequired: false,
	},
	{
		URI:          "/posts/{id}",
		Method:       http.MethodGet,
		Function:     controllers.FindPost,
		AuthRequired: false,
	},
	{
		URI:          "/posts/{id}",
		Method:       http.MethodPut,
		Function:     controllers.UpdatePost,
		AuthRequired: false,
	},
	{
		URI:          "/posts/{id}",
		Method:       http.MethodDelete,
		Function:     controllers.DeletePost,
		AuthRequired: false,
	},
	{
		URI:          "/users/{id}/posts",
		Method:       http.MethodGet,
		Function:     controllers.UserPosts,
		AuthRequired: false,
	},
	{
		URI:          "/posts/{id}/like",
		Method:       http.MethodPost,
		Function:     controllers.LikePost,
		AuthRequired: false,
	},
	{
		URI:          "/posts/{id}/unlike",
		Method:       http.MethodPost,
		Function:     controllers.UnlikePost,
		AuthRequired: false,
	},
	{
		URI:          "/posts/{id}/likes",
		Method:       http.MethodGet,
		Function:     controllers.LikesPost,
		AuthRequired: false,
	},
	{
		URI:          "/posts/{id}/comments",
		Method:       http.MethodPost,
		Function:     controllers.NewComment,
		AuthRequired: false,
	},
	{
		URI:          "/posts/{id}/comments",
		Method:       http.MethodDelete,
		Function:     controllers.DeleteComment,
		AuthRequired: false,
	},
}
