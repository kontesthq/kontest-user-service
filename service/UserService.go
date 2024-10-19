package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"kontest-user-service/database"
	error2 "kontest-user-service/error"
	"kontest-user-service/model"
	"kontest-user-service/utils"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) GetUser(uid uuid.UUID) (*model.GetUserResponse, error) {
	user, err := database.FindUserByID(uid)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, &error2.UserNotFoundError{}
	}

	email, err := getEmail(uid.String())

	if err != nil {
		slog.Error(fmt.Sprintf("Error in getting email: %v\n", err))
	}

	return &model.GetUserResponse{
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		Email:               email,
		LeetcodeUsername:    user.LeetcodeUsername,
		CodechefUsername:    user.CodechefUsername,
		CodeforcesUsername:  user.CodeforcesUsername,
		Sites:               user.Sites,
		MinDurationInSecond: user.MinDurationInSecond,
		MaxDurationInSecond: user.MaxDurationInSecond,
	}, nil
}

func getEmail(userID string) (string, error) {
	apiGatewayURL := "http" + "://" + utils.API_GATEWAY_HOST + ":" + strconv.Itoa(utils.API_GATEWAY_PORT)
	slog.Info(fmt.Sprintf("API Gateway URL: %s\n", apiGatewayURL))
	urlString := apiGatewayURL + "/auth/internal/email"
	slog.Info(fmt.Sprintf("FINAL URL: %s\n", urlString))

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("error parsing base URL: %w", err)
	}

	query := url.Values{}
	query.Add("user_id", userID)
	parsedURL.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating request: %v\n", err))
		return "", err
	}

	slog.Info(fmt.Sprintf("Final Full URL: %s\n", parsedURL.String()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error(fmt.Sprintf("Error making request: %v\n", err))
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error(fmt.Sprintf("Error: received status code %d\n", resp.StatusCode))
		return "", err
	}

	email, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("Error reading response body: %v\n", err))
		return "", err
	}

	return strings.Trim(strings.TrimSpace(string(email)), "\""), nil
}

func getUser(uid uuid.UUID) (*model.User, error) {
	currentSavedUser, err := database.FindUserByID(uid)

	if err != nil {
		if errors.Is(err, &error2.UserNotFoundError{}) {
			// Create a new user with the provided UID
			newUser := &model.User{
				ID:                  uid,
				FirstName:           "",
				LastName:            "",
				AccountCreateDate:   time.Now(),
				LeetcodeUsername:    "",
				CodechefUsername:    "",
				CodeforcesUsername:  "",
				Sites:               []model.Site{},
				MinDurationInSecond: 0,
				MaxDurationInSecond: 0,
			}

			// create the user
			wasOperationSuccessful, err := database.UpdateUserOrCreate(newUser, nil)
			if err != nil {
				// Return error if the creation failed
				return nil, err
			}

			if !wasOperationSuccessful {
				// Return an error if the upsert operation was unsuccessful for some reason
				return nil, errors.New("failed to create user")
			}

			return newUser, nil

		} else {
			slog.Error("Error in updating user: ", slog.String("error", err.Error()))
			return nil, errors.New("internal server error")
		}
	}

	return currentSavedUser, nil
}

func (us *UserService) UpdateUser(uid uuid.UUID, request model.PutUserRequest) (bool, error) {
	currentSavedUser, err := getUser(uid)

	if err != nil {
		return false, err
	}

	fmt.Printf("currentSavedUser: %v\n", currentSavedUser)

	if request.FirstName != nil {
		currentSavedUser.FirstName = *request.FirstName
	}

	if request.LastName != nil {
		currentSavedUser.LastName = *request.LastName
	}

	if request.LeetcodeUsername != nil {
		currentSavedUser.LeetcodeUsername = *request.LeetcodeUsername
	}

	if request.CodechefUsername != nil {
		currentSavedUser.CodechefUsername = *request.CodechefUsername
	}

	if request.CodeforcesUsername != nil {
		currentSavedUser.CodeforcesUsername = *request.CodeforcesUsername
	}

	if request.Sites != nil {
		currentSavedUser.Sites = request.Sites
	}

	if request.MinDurationInSecond != nil {
		currentSavedUser.MinDurationInSecond = *request.MinDurationInSecond
	}

	if request.MaxDurationInSecond != nil {
		currentSavedUser.MaxDurationInSecond = *request.MaxDurationInSecond
	}

	isSuccessful, err := database.UpdateUserOrCreate(currentSavedUser, nil)
	if err != nil {
		return false, err
	}

	if !isSuccessful {
		return false, nil

	}

	return true, nil
}
