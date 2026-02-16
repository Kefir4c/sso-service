package models

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	PassHash []byte `json:"-"`
}

func NewUser(id int64, email string, passhash []byte) *User {
	return &User{
		ID:       id,
		Email:    email,
		PassHash: passhash,
	}
}
