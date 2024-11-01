package routes

import (
	"kontest-user-service/handler"
	"kontest-user-service/handler/admin"
	"net/http"
)

func RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /user_info", handler.GetUserHandler)
	router.HandleFunc("PUT /user_info", handler.PutUserHandler)
}

func RegisterAdminRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /admin/user_info", admin.GetUserHandler)
}
