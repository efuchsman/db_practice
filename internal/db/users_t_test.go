package db

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testUserEli = &User{
		Id:          "12infioed",
		FirstName:   "Eli",
		LastName:    "Fuchsman",
		Email:       "testEmail@mail.com",
		Address:     "1123 Street St.",
		City:        "Denver",
		State:       "CO",
		ZipCode:     "80108",
		DateOfBirth: "12/14/1993",
	}

	testUserEli2 = &User{
		Id:          "12infio8ed",
		FirstName:   "Eli",
		LastName:    "Fuchsman",
		Email:       "testEmail2@mail.com",
		Address:     "1123 Street St.",
		City:        "Denver",
		State:       "CO",
		ZipCode:     "80108",
		DateOfBirth: "12/14/1993",
	}
)

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		description string
		testUser    *User
		expectedErr error
	}{
		{
			description: "Success: User added to the DB",
			testUser:    testUserEli,
			expectedErr: nil,
		},
		{
			description: "Failure: Email already exists",
			testUser:    testUserEli2,
			expectedErr: ErrEmailExists,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			db, err := NewDB(connStr, true, t.Name())
			// Use txdb for testing
			require.NoError(t, err)
			defer db.Close()

			tx, err := db.Conn.Begin()
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback()

			_, err = db.CreateUser(testUserEli2.FirstName, testUserEli2.LastName, testUserEli2.Email, testUserEli2.Address, testUserEli2.City, testUserEli2.State, testUserEli2.ZipCode, testUserEli2.DateOfBirth)

			user, err := db.CreateUser(tc.testUser.FirstName, tc.testUser.LastName, tc.testUser.Email, tc.testUser.Address, tc.testUser.City, tc.testUser.State, tc.testUser.ZipCode, tc.testUser.DateOfBirth)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NotNil(t, user)
				assert.NotNil(t, user.Id)
				assert.NoError(t, err, tc.description)
				require.NoError(t, err)
				assert.Equal(t, tc.testUser.FirstName, user.FirstName)
				assert.Equal(t, tc.testUser.LastName, user.LastName)
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	testCases := []struct {
		description string
		email       string
		expectedErr error
	}{
		{
			description: "Success: User added to the DB",
			email:       testUserEli.Email,
			expectedErr: nil,
		},
		{
			description: "Failure: No user found",
			email:       testUserEli2.Email,
			expectedErr: sql.ErrNoRows,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			db, err := NewDB(connStr, true, t.Name())
			// Use txdb for testing
			require.NoError(t, err)
			defer db.Close()

			tx, err := db.Conn.Begin()
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback()

			_, err = db.CreateUser(testUserEli.FirstName, testUserEli.LastName, testUserEli.Email, testUserEli.Address, testUserEli.City, testUserEli.State, testUserEli.ZipCode, testUserEli.DateOfBirth)
			require.NoError(t, err)

			foundUser, err := db.GetUserByEmail(tc.email)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NotNil(t, foundUser)
				assert.NotNil(t, foundUser.Id)
				assert.NoError(t, err, tc.description)
				require.NoError(t, err)
				assert.Equal(t, testUserEli.FirstName, foundUser.FirstName)
				assert.Equal(t, testUserEli.LastName, foundUser.LastName)
			}
		})
	}
}

func TestGetUserById(t *testing.T) {
	testCases := []struct {
		description string
		id          string
		expectedErr error
	}{
		{
			description: "Success: User added to the DB",
			id:          testUserEli.Id,
			expectedErr: nil,
		},
		{
			description: "Failure: No user found",
			id:          testUserEli2.Id,
			expectedErr: sql.ErrNoRows,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			db, err := NewDB(connStr, true, t.Name())
			// Use txdb for testing
			require.NoError(t, err)
			defer db.Close()

			tx, err := db.Conn.Begin()
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback()

			_, err = db.CreateUser(testUserEli.FirstName, testUserEli.LastName, testUserEli.Email, testUserEli.Address, testUserEli.City, testUserEli.State, testUserEli.ZipCode, testUserEli.DateOfBirth)
			require.NoError(t, err)

			foundUser, err := db.GetUserById(tc.id)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NotNil(t, foundUser)
				assert.NotNil(t, foundUser.Id)
				assert.NoError(t, err, tc.description)
				require.NoError(t, err)
				assert.Equal(t, testUserEli.FirstName, foundUser.FirstName)
				assert.Equal(t, testUserEli.LastName, foundUser.LastName)
			}
		})
	}
}
