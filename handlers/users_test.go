package handlers

import (
	"bytes"
	"db_practice/internal/users"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testUserEli = &users.User{
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

	testUserEli2 = &users.User{
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

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		description  string
		userClient   *users.TestClient
		url          string
		requestBody  *CreateUserRequest
		expectedBody string
		expectedCode int
	}{
		{
			description: "Success: User created",
			userClient: &users.TestClient{
				CreateUserData: testUserEli,
			},
			url: "/create",
			requestBody: &CreateUserRequest{
				FirstName:   "Eli",
				LastName:    "Fuchsman",
				Email:       "testEmail@mail.com",
				Address:     "1123 Street St.",
				City:        "Denver",
				State:       "CO",
				ZipCode:     "80108",
				DateOfBirth: "12/14/1993",
			},
			expectedBody: `{"id":"12infioed","first_name":"Eli","last_name":"Fuchsman","email":"testEmail@mail.com","address":"1123 Street St.","city":"Denver","state":"CO","zip":"80108","dob":"12/14/1993"}`,
			expectedCode: 201,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.description)

			h := NewUsersHandler(tc.userClient)
			body, _ := json.Marshal(tc.requestBody)
			r := httptest.NewRequest("POST", tc.url, bytes.NewBuffer(body))

			w := httptest.NewRecorder()
			h.CreateUser(w, r)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}
