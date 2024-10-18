package handler

import (
	"encoding/json"
	"kontest-user-service/model"
	"net/http"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	getUserRequest := model.GetUserRequest{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&getUserRequest); err != nil {
		http.Error(w, "Please provide login request in correct JSON format", http.StatusBadRequest)
		return
	}

}
