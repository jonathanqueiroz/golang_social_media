package routes

import (
	"database/sql"
	"net/http"
	"project01/src/controllers"
)

func profileRoutes(db *sql.DB) []Route {
	profileController := controllers.NewProfileController(db)

	return []Route{
		{
			URI:          "/profile",
			Method:       http.MethodGet,
			Function:     profileController.GetProfile,
			AuthRequired: true,
		},
	}
}
