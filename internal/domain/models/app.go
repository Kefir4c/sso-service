package models

// App represents registered application that can use SSO service.
// Each app has its own secret key for JWT signing.
type App struct {
	ID     int
	Name   string
	Secret string
}
