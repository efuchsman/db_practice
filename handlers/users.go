package handlers

import (
	"db_practice/internal/users"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
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
	ZipCode     string `json:"zip"`
	DateOfBirth string `json:"dob"`
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

	req.Email = strings.ToLower(req.Email)
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
		return
	}

	user, err := u.usersClient.CreateUser(req.FirstName, req.LastName, req.Email, req.Address, req.City, req.State, req.ZipCode, req.DateOfBirth)
	if err != nil {
		log.WithFields(loggedFields).Errorf("%+v", err)
		InternalError500(w, "Users", err)
		return
	}

	Created201(w, user)
}

func (u *UsersHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /users/:email request")

	vars := mux.Vars(r)
	email, exists := vars["email"]
	email = strings.ToLower(email)
	fields := log.Fields{"Email": email}
	if !exists || email == "" {
		log.WithFields(fields).Error("MISSING_ARG_EMAIL")
		BadRequest400(w, "Users", "MISSING_ARG_EMAIL")
		return
	}

	user, err := u.usersClient.GetUserByEmail(email)
	if err != nil {
		log.WithFields(fields).Errorf("%+v", err)
		NotFound404(w, "Users")
		return
	}
	OK200(w, user)
}
