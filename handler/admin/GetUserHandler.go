package admin

import (
	"encoding/json"
	"kontest-user-service/model"
	"kontest-user-service/service"
	"net/http"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	getUserRequest := model.GetUserRequest{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&getUserRequest); err != nil {
		http.Error(w, "Please provide request body in correct JSON format", http.StatusBadRequest)
		return
	}

	userService := service.NewUserService()

	user, err := userService.GetUser(getUserRequest.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&user)
}
