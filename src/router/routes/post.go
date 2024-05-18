package routes

import (
	"database/sql"
	"net/http"
	"project01/src/controllers"
)

func postRoutes(db *sql.DB) []Route {
	postController := controllers.NewPostController(db)

	return []Route{
		{
			URI:          "/posts",
			Method:       http.MethodPost,
			Function:     postController.NewPost,
			AuthRequired: true,
		},
		{
			URI:          "/posts",
			Method:       http.MethodGet,
			Function:     postController.PostsFollowedUsers,
			AuthRequired: true,
		},
		{
			URI:          "/posts/{id}",
			Method:       http.MethodGet,
			Function:     postController.FindPost,
			AuthRequired: true,
		},
		{
			URI:          "/posts/{id}",
			Method:       http.MethodPut,
			Function:     postController.UpdatePost,
			AuthRequired: true,
		},
		{
			URI:          "/posts/{id}",
			Method:       http.MethodDelete,
			Function:     postController.DeletePost,
			AuthRequired: true,
		},
		{
			URI:          "/users/{id}/posts",
			Method:       http.MethodGet,
			Function:     postController.UserPosts,
			AuthRequired: true,
		},
		{
			URI:          "/posts/{id}/like",
			Method:       http.MethodPost,
			Function:     postController.LikePost,
			AuthRequired: true,
		},
		{
			URI:          "/posts/{id}/unlike",
			Method:       http.MethodPost,
			Function:     postController.UnlikePost,
			AuthRequired: true,
		},
		{
			URI:          "/posts/{id}/likes",
			Method:       http.MethodGet,
			Function:     postController.LikesPost,
			AuthRequired: true,
		},
		{
			URI:          "/posts/{id}/comments",
			Method:       http.MethodPost,
			Function:     postController.NewPost,
			AuthRequired: true,
		},
	}
}
