package db

import (
	"database/sql"
	"fmt"
)

type Client interface {
	CreateUser(firstName, lastName, email, address, city, state, zip, dob string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserById(id string) (*User, error)
}

type DB struct {
	Conn *sql.DB
}

func NewDB(connStr string) (*DB, error) {
	var db *sql.DB
	var err error

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the database")
	return &DB{
		Conn: db,
	}, nil
}

func (db *DB) Close() {
	db.Conn.Close()
	fmt.Println("Closed the database connection")
}
