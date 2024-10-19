package handler

import (
	"encoding/json"
	"fmt"
	"kontest-user-service/model"
	"kontest-user-service/service"
	"kontest-user-service/utils"
	"log/slog"
	"net/http"
)

func PutUserHandler(w http.ResponseWriter, r *http.Request) {
	putUserRequest := model.PutUserRequest{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&putUserRequest); err != nil {
		http.Error(w, "Please provide request body in correct JSON format", http.StatusBadRequest)
		return
	}

	userId := r.Header.Get(utils.UserIdRequestHeader)

	if userId == "" {
		http.Error(w, "JWT token not provided", http.StatusBadRequest)
		return
	}

	uid, err := utils.IsValidUUID(userId)

	if err != nil {
		slog.Error(fmt.Sprintf("userId: %s not valid", userId))
		http.Error(w, "JWT token not valid", http.StatusBadRequest)
		return
	}

	userService := service.NewUserService()

	hasUserUpdated, err := userService.UpdateUser(uid, putUserRequest)
	if err != nil {
		slog.Error(fmt.Sprintf("Error in updating user: %s with error: %v", userId, err))
		http.Error(w, "Error in updating user", http.StatusInternalServerError)
		return
	}

	if !hasUserUpdated {
		slog.Error(fmt.Sprintf("Error in updating user: %s without any error", userId))
		http.Error(w, "Error in updating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated"))
}
