package routes

import (
	"database/sql"
	"net/http"
	"project01/src/controllers"
)

func profileRoutes(db *sql.DB) []Route {
	postController := controllers.NewPostController(db)
	userController := controllers.NewUserController(db)
	profileController := controllers.NewProfileController(userController.GetRepo(), postController.GetRepo())

	return []Route{
		{
			URI:          "/profile",
			Method:       http.MethodGet,
			Function:     profileController.GetProfile,
			AuthRequired: true,
		},
	}
}
