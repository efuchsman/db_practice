package users

import (
	"database/sql"
	"db_practice/internal/db"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testUserEli = &db.User{
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

	testUserEli2 = &db.User{
		Id:          "12infioEd",
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

var connStr string

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Error getting the current file path.")
	}

	dir := filepath.Dir(filename)
	projectRoot := filepath.Join(dir, "..", "..")
	configPath := filepath.Join(projectRoot, "config", "test.yml")
	if err := godotenv.Load(filepath.Join(projectRoot, ".env")); err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic("Error reading testing configuration file: " + err.Error())
	}
	viper.SetDefault("environment.test.database.user", os.Getenv("TEST_USER"))
	viper.SetDefault("environment.test.database.password", os.Getenv("TEST_PASSWORD"))
	viper.SetDefault("environment.test.database.name", os.Getenv("TEST_DB"))
	viper.SetDefault("environment.test.database.connection_string", os.Getenv("TEST_CONN_STR"))

	connStr = viper.GetString("environment.test.database.connection_string")
	if connStr == "" {
		panic("Connection string not found in configuration")
	}
}

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		description    string
		db             *db.TestClient
		expectedOutput *User
		expectedErr    error
	}{
		{
			description: "Success: User created",
			db: &db.TestClient{
				CreateUserData: testUserEli,
			},
			expectedOutput: &User{
				Id:          "12infioed",
				FirstName:   "Eli",
				LastName:    "Fuchsman",
				Email:       "testEmail@mail.com",
				Address:     "1123 Street St.",
				City:        "Denver",
				State:       "CO",
				ZipCode:     "80108",
				DateOfBirth: "12/14/1993",
			},
			expectedErr: nil,
		},
		{
			description: "Failure: Email already exists",
			db: &db.TestClient{
				CreateUserErr: db.ErrEmailExists,
			},
			expectedErr: db.ErrEmailExists,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			c := NewUsersClient(tc.db)
			user, err := c.CreateUser(testUserEli.FirstName, testUserEli.LastName, testUserEli.Email, testUserEli.Address, testUserEli.City, testUserEli.State, testUserEli.ZipCode, testUserEli.DateOfBirth)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NotNil(t, user)
				assert.NotNil(t, user.Id)
				assert.NoError(t, err, tc.description)
				require.NoError(t, err)
				assert.Equal(t, tc.expectedOutput.FirstName, user.FirstName)
				assert.Equal(t, tc.expectedOutput.LastName, user.LastName)
			}
		})
	}
}

func TestGetUserById(t *testing.T) {
	testCases := []struct {
		description    string
		id             string
		db             *db.TestClient
		expectedOutput *User
		expectedErr    error
	}{
		{
			description: "Success: User found",
			id:          testUserEli.Id,
			db: &db.TestClient{
				CreateUserData:  testUserEli,
				GetUserByIdData: testUserEli,
			},
			expectedOutput: &User{
				Id:          "12infioed",
				FirstName:   "Eli",
				LastName:    "Fuchsman",
				Email:       "testEmail@mail.com",
				Address:     "1123 Street St.",
				City:        "Denver",
				State:       "CO",
				ZipCode:     "80108",
				DateOfBirth: "12/14/1993",
			},
			expectedErr: nil,
		},
		{
			description: "Failure: No user found",
			id:          testUserEli.Id,
			db: &db.TestClient{
				CreateUserData: testUserEli,
				GetUserByIdErr: sql.ErrNoRows,
			},
			expectedOutput: nil,

			expectedErr: sql.ErrNoRows,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			c := NewUsersClient(tc.db)

			_, err := c.CreateUser(testUserEli.FirstName, testUserEli.LastName, testUserEli.Email, testUserEli.Address, testUserEli.City, testUserEli.State, testUserEli.ZipCode, testUserEli.DateOfBirth)
			require.NoError(t, err)

			user, err := c.GetUserById(tc.id)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NotNil(t, user)
				assert.NoError(t, err, tc.description)
				assert.Equal(t, user, tc.expectedOutput)
			}
		})
	}
}
