package db

type TestClient struct {
	CreateUserData *User
	CreateUserErr  error

	GetUserByEmailData *User
	GetUserByEmailErr  error

	GetUserByIdData *User
	GetUserByIdErr  error
}

func (c TestClient) CreateUser(firstName, lastName, email, address, city, state, zip, dob string) (*User, error) {
	return c.CreateUserData, c.CreateUserErr
}

func (c TestClient) GetUserByEmail(email string) (*User, error) {
	return c.GetUserByEmailData, c.GetUserByEmailErr
}

func (c TestClient) GetUserById(id string) (*User, error) {
	return c.GetUserByIdData, c.GetUserByIdErr
}
