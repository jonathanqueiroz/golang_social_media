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
		URI:          "/posts",
		Method:       http.MethodGet,
		Function:     controllers.AllPosts,
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
}
