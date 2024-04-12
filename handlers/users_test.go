package handlers

import (
	"db_practice/internal/users"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var (
	testUserEli = &users.User{
		Id:          "12infioed",
		FirstName:   "Eli",
		LastName:    "Fuchsman",
		Email:       "testemail@mail.com",
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
		Email:       "testemail2@mail.com",
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
		requestBody  io.Reader
		expectedBody string
		expectedCode int
	}{
		{
			description: "Success: User created",
			userClient: &users.TestClient{
				CreateUserData: testUserEli,
			},
			url: "/create",
			requestBody: strings.NewReader(`{
				"first_name": "Eli",
				"last_name": "Fuchsman",
				"email": "testemail@mail.com",
				"address": "1123 Street St.",
				"city": "Denver",
				"state": "CO",
				"zip": "80108",
				"dob": "12/14/1993"
			}`),
			expectedBody: `{"id":"12infioed","first_name":"Eli","last_name":"Fuchsman","email":"testemail@mail.com","address":"1123 Street St.","city":"Denver","state":"CO","zip":"80108","dob":"12/14/1993"}`,
			expectedCode: 201,
		},
		{
			description: "Failure: Missing field",
			url:         "/create",
			requestBody: strings.NewReader(`{
					"first_name": "",
					"last_name": "Fuchsman",
					"email": "testemail@mail.com",
					"address": "1123 Street St.",
					"city": "Denver",
					"state": "CO",
					"zip": "80108",
					"dob": "12/14/1993"
			}`),
			expectedBody: `{"message":"BAD_REQUEST","resource":"Users","description":"The value provided is invalid.","errors":[{"field":"First Name","error_code":"invalid"}]}`,
			expectedCode: 400,
		},
		{
			description: "Failure: Bad JSON",
			url:         "/create",
			requestBody: strings.NewReader(`{
				"first_name": "Eli"
				"last_name": "Fuchsman",
				"email": "testemail@mail.com",
				"address": "1123 Street St.",
				"city": "Denver",
				"state": "CO",
				"zip": "80108",
				"dob": "12/14/1993"
			}`),
			expectedBody: `{"message":"BAD_REQUEST","resource":"Users","description":"The value provided is invalid.","errors":[{"field":"INVALID_JSON","error_code":"invalid"}]}`,
			expectedCode: 400,
		},
		{
			description: "Failure: Internal Error",
			url:         "/create",
			userClient: &users.TestClient{
				CreateUserErr: errors.New("error"),
			},
			requestBody: strings.NewReader(`{
				"first_name": "Eli",
				"last_name": "Fuchsman",
				"email": "testemail@mail.com",
				"address": "1123 Street St.",
				"city": "Denver",
				"state": "CO",
				"zip": "80108",
				"dob": "12/14/1993"
			}`),
			expectedBody: `{"message":"INTERNAL_ERROR","resource":"Users","description":"An internal error occurred."}`,
			expectedCode: 500,
		},
		{
			description: "Failure: Email in use",
			url:         "/create",
			userClient: &users.TestClient{
				GetUserByEmailData: testUserEli,
			},
			requestBody: strings.NewReader(`{
				"first_name": "Eli",
				"last_name": "Fuchsman",
				"email": "testemail@mail.com",
				"address": "1123 Street St.",
				"city": "Denver",
				"state": "CO",
				"zip": "80108",
				"dob": "12/14/1993"
			}`),
			expectedBody: `{"message":"CONFLICT_ERROR","resource":"Users","description":"there is a conflict with your request"}`,
			expectedCode: 409,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Log(tc.description)

			h := NewUsersHandler(tc.userClient)
			r := httptest.NewRequest("POST", tc.url, tc.requestBody)

			w := httptest.NewRecorder()
			h.CreateUser(w, r)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	testCases := []struct {
		description  string
		userClient   *users.TestClient
		email        string
		expectedBody string
		expectedCode int
	}{
		{
			description: "Success: User found",
			userClient: &users.TestClient{
				GetUserByEmailData: testUserEli,
			},
			email:        testUserEli.Email,
			expectedBody: `{"id":"12infioed","first_name":"Eli","last_name":"Fuchsman","email":"testemail@mail.com","address":"1123 Street St.","city":"Denver","state":"CO","zip":"80108","dob":"12/14/1993"}`,
			expectedCode: 200,
		},
		{
			description:  "Failure: No Email",
			email:        "",
			expectedBody: `{"message":"BAD_REQUEST","resource":"Users","description":"The value provided is invalid.","errors":[{"field":"MISSING_ARG_EMAIL","error_code":"invalid"}]}`,
			expectedCode: 400,
		},
		{
			description: "Failure: No User",
			userClient: &users.TestClient{
				GetUserByEmailErr: errors.New("NOT_FOUND"),
			},
			email:        testUserEli2.Email,
			expectedBody: `{"message":"NOT_FOUND","resource":"Users","description":"What you are looking for cannot be found."}`,
			expectedCode: 404,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Log(tc.description)

			h := NewUsersHandler(tc.userClient)
			url := fmt.Sprintf("/users/%s", tc.email)
			r := httptest.NewRequest("GET", url, nil)
			r = mux.SetURLVars(r, map[string]string{"email": tc.email})

			w := httptest.NewRecorder()
			h.GetUserByEmail(w, r)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}
