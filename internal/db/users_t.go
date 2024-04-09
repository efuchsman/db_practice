package db

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Table "public.users"
// Column   |          Type          | Collation | Nullable | Default
// ------------+------------------------+-----------+----------+---------
// id         | character varying(10)  |           | not null |
// first_name | character varying(50)  |           | not null |
// last_name  | character varying(50)  |           | not null |
// email      | character varying(100) |           | not null |
// address    | character varying(255) |           |          |
// city       | character varying(100) |           |          |
// state      | character varying(100) |           |          |
// zip        | character varying(20)  |           |          |
// dob        | character varying(20)  |           |          |
// Indexes:
//
//	"users_pkey" PRIMARY KEY, btree (id)
//	"users_email_key" UNIQUE CONSTRAINT, btree (email)

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

var ErrEmailExists = errors.New("Email is already in use")
var ErrIdExists = errors.New("Unique id required")

func (db *DB) CreateUser(firstName, lastName, email, address, city, state, zip, dob string) (*User, error) {
	existingEmail, err := db.GetUserByEmail(email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingEmail != nil {
		return nil, ErrEmailExists
	}

	id, err := generateUserId(db)
	if err != nil {
		return nil, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		log.Errorf("failed to begin transaction: %v", err)
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != sql.ErrTxDone && err != nil {
			log.Errorf("failed to rollback transaction: %v", err)
		}
	}()

	newUser := &User{
		Id:          id,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		Address:     address,
		City:        city,
		State:       state,
		ZipCode:     zip,
		DateOfBirth: dob,
	}

	createdUser, err := db.CreateUserTx(tx, newUser)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("failed to commit transaction: %v", err)
		return nil, err
	}
	return createdUser, nil
}

func (db *DB) CreateUserTx(tx *sql.Tx, u *User) (*User, error) {

	query := `
			INSERT INTO Users (id, first_name, last_name, email, address, city, state, zip, dob)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (email) DO NOTHING
			RETURNING id, first_name, last_name, email, address, city, state, zip, dob;`

	var id, firstName, lastName, email, address, city, state, zip, dob string
	err := tx.QueryRow(query, u.Id, u.FirstName, u.LastName, u.Email, u.Address, u.City, u.State, u.ZipCode, u.DateOfBirth).Scan(&id, &firstName, &lastName, &email, &address, &city, &state, &zip, &dob)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user already exists with email %s", email)
		}
		return nil, err
	}

	newUser := &User{
		Id:          u.Id,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Email:       u.Email,
		Address:     u.Address,
		City:        u.City,
		State:       u.State,
		ZipCode:     u.ZipCode,
		DateOfBirth: u.DateOfBirth,
	}

	fmt.Printf("User inserted successfully with ID: %s\n", id)
	return newUser, nil
}

func (db *DB) GetUserByEmail(email string) (*User, error) {

	query := `
		SELECT id, first_name, last_name, email, address, city, state, zip, dob
		FROM users
		WHERE email = $1
  `

	row := db.Conn.QueryRow(query, email)

	var user User
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Address, &user.City, &user.State, &user.ZipCode, &user.DateOfBirth)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}

		return nil, err
	}

	return &user, nil
}

func (db *DB) GetUserById(id string) (*User, error) {

	query := `
		SELECT id, first_name, last_name, email, address, city, state, zip, dob
		FROM users
		WHERE id = $1
  `

	row := db.Conn.QueryRow(query, id)

	var user User
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Address, &user.City, &user.State, &user.ZipCode, &user.DateOfBirth)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}

		return nil, err
	}

	return &user, nil
}

func generateUserId(db *DB) (string, error) {
	for {
		input := make([]byte, 16)
		if _, err := rand.Read(input); err != nil {
			return "", err
		}

		hash := sha256.Sum256(input)
		id := hex.EncodeToString(hash[:])[:10]

		existingId, err := db.GetUserById(id)
		if err != nil && err != sql.ErrNoRows {
			return "", err
		}
		if existingId == nil {
			return id, nil
		}
	}
}
