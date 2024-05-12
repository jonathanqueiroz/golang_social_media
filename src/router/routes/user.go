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
		Function:     controllers.AllUsers,
		AuthRequired: true,
	},
	{
		URI:          "/users/find-by-name",
		Method:       http.MethodGet,
		Function:     controllers.FindUserByName,
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
}
