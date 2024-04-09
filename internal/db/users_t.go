package db

import (
	"database/sql"
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

func (db *DB) CreateUser(firstName, lastName, email, address, city, state, zip, dob string) (*User, error) {

	return nil, nil
}

func (db *DB) CreateUserTx(tx *sql.Tx, u *User) (*User, error) {

	query := `
			INSERT INTO Users (id, first_name, last_name, email, address, city, state, zip, dob)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (id) DO NOTHING
			ON CONFLICT (email) DO NOTHING
			RETURNING id, first_name, last_name, email, address, city, state, zip, dob;`

	var userID int
	var id, firstName, lastName, email, address, city, state, zip, dob string
	err := tx.QueryRow(query, u.Id, u.FirstName, u.LastName, u.Email, u.Address, u.City, u.State, u.ZipCode, u.DateOfBirth).Scan(&id, &firstName, &lastName, &email, &address, &city, &state, &zip, &dob)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("failed to insert transactional user: %v", err)
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

	fmt.Printf("User inserted successfully with ID: %d\n", userID)
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
			return nil, fmt.Errorf("user with email %s not found", email)
		}

		return nil, err
	}

	return &user, nil
}
