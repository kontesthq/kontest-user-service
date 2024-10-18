package database

import (
	"fmt"
	"github.com/google/uuid"
	error2 "kontest-user-service/error"
	"kontest-user-service/model"
)

func FindUserByID(userID uuid.UUID) (*model.User, error) {
	var user model.User

	query := `SELECT * FROM user_info WHERE id = :id`

	// Use NamedQuery to retrieve the user by ID
	rows, err := GetDB().NamedQuery(query, map[string]interface{}{
		"id": userID,
	})

	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	// Check if any rows are returned
	if rows.Next() {
		if err := rows.StructScan(&user); err != nil {
			return nil, fmt.Errorf("error scanning result: %v", err)
		}
		return &user, nil
	}

	// If no rows were returned, the user was not found
	return nil, &error2.UserNotFoundError{}
}
