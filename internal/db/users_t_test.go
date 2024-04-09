package db

import (
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
