package models

// User represents system user with authentication data.
type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	PassHash []byte `json:"pass_hash"`
}
