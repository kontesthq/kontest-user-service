package handler

import (
	"encoding/json"
	"kontest-user-service/service"
	"kontest-user-service/utils"
	"log/slog"
	"net/http"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	loggedInUserID := r.Header.Get(utils.UserIdRequestHeader)

	if loggedInUserID == "" {
		slog.Error(utils.UserIdRequestHeader + " header not present")
		http.Error(w, "Internal error occurred", http.StatusInternalServerError)
		return
	}

	uid, err := utils.IsValidUUID(loggedInUserID)

	if err != nil {
		slog.Error("Invalid user id")
		http.Error(w, "Internal error occurred", http.StatusInternalServerError)
		return
	}

	userService := service.NewUserService()

	user, err := userService.GetUser(uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&user)
}
