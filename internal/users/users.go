package users

import (
	"db_practice/internal/db"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Client interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(id string) (*User, error)
	CreateUser(firstName, lastName, email, address, city, state, zip, dob string) (*User, error)
}

type UsersClient struct {
	db db.Client
}

type User struct {
	Id          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
	ZipCode     string `json:"zip"`
	DateOfBirth string `json:"dob"`
}

func NewUsersClient(data db.Client) *UsersClient {
	return &UsersClient{
		db: data,
	}
}

func (u *UsersClient) GetUserByEmail(email string) (*User, error) {
	fields := log.Fields{"Email": email}

	user, err := u.db.GetUserByEmail(email)
	if err != nil {
		log.WithFields(fields).Errorf("User not found with email: %s", email)
		return nil, errors.WithStack(err)
	}

	foundUser := &User{
		Id:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Address:     user.Address,
		City:        user.City,
		State:       user.State,
		ZipCode:     user.ZipCode,
		DateOfBirth: user.DateOfBirth,
	}

	return foundUser, nil
}

func (u *UsersClient) GetUserById(id string) (*User, error) {
	fields := log.Fields{"Id": id}

	user, err := u.db.GetUserById(id)
	if err != nil {
		log.WithFields(fields).Errorf("User not found with id: %s", id)
		return nil, errors.WithStack(err)
	}

	foundUser := &User{
		Id:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Address:     user.Address,
		City:        user.City,
		State:       user.State,
		ZipCode:     user.ZipCode,
		DateOfBirth: user.DateOfBirth,
	}

	return foundUser, nil
}

func (u *UsersClient) CreateUser(firstName, lastName, email, address, city, state, zip, dob string) (*User, error) {
	fields := log.Fields{"First Name": firstName, "Last Name": lastName, "Email": email, "Address": address, "City": city, "State": state, "Zip Code": zip, "Date of Birth": dob}

	user, err := u.db.CreateUser(firstName, lastName, email, address, city, state, zip, dob)
	if err != nil {
		log.WithFields(fields).Errorf("Failed to create user: %+v", err)
		return nil, errors.WithStack(err)
	}

	newUser := &User{
		Id:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Address:     user.Address,
		City:        user.City,
		State:       user.State,
		ZipCode:     user.ZipCode,
		DateOfBirth: user.DateOfBirth,
	}

	return newUser, nil
}
