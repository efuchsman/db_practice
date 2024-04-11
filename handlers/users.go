package handlers

import (
	"db_practice/internal/users"
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type UsersHandler struct {
	usersClient users.Client
}

type CreateUserRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
	ZipCode     string `json:"zip_code"`
	DateOfBirth string `json:"date_of_birth"`
}

func NewUsersHandler(u users.Client) *UsersHandler {
	return &UsersHandler{
		usersClient: u,
	}
}

func (u *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest400(w, "Users", "INVALID_JSON")
		return
	}
	defer r.Body.Close()

	fields := map[string]string{"First Name": req.FirstName, "Last Name": req.LastName, "Email": req.Email, "Address": req.Address, "City": req.City, "State": req.State, "Zip Code": req.ZipCode, "Date of Birth": req.DateOfBirth}

	missingFields := []string{}

	for field, value := range fields {
		if value == "" {
			missingFields = append(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		fields := log.Fields{"missing_fields": strings.Join(missingFields, ", ")}
		log.WithFields(fields).Error("Missing required fields")
	}

	if len(missingFields) > 0 {
		BadRequest400(w, "Users", strings.Join(missingFields, ", "))
		return
	}

	loggedFields := log.Fields{"First Name": req.FirstName, "Last Name": req.LastName, "Email": req.Email, "Address": req.Address, "City": req.City, "State": req.State, "Zip Code": req.ZipCode, "Date of Birth": req.DateOfBirth}

	emailCheck, _ := u.usersClient.GetUserByEmail(fields["Email"])
	if emailCheck != nil {
		log.WithFields(loggedFields).Errorf("Email already in use: %s", req.Email)
		ConflictError409(w, "Users", "email")
	}

	user, err := u.usersClient.CreateUser(req.FirstName, req.LastName, req.Email, req.Address, req.City, req.State, req.ZipCode, req.DateOfBirth)
	if err != nil {
		log.WithFields(loggedFields).Errorf("%+v", err)
		InternalError500(w, "Users", err)
	}

	Created201(w, user)
}
