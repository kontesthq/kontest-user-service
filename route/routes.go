package routes

import (
	"kontest-user-service/handler"
	"net/http"
)

func RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /userInfo", handler.GetUserHandler)
	router.HandleFunc("PUT /userInfo", handler.PutUserHandler)
}
