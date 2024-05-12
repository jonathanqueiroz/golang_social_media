package routes

import (
	"net/http"
	"project01/src/controllers"
)

var loginRoutes = []Route{
	{
		URI:          "/login",
		Method:       http.MethodPost,
		Function:     controllers.Login,
		AuthRequired: false,
	},
}
