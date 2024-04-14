package entity

type User struct {
	ID           int
	Username     string
	PasswordHash string
	Role         string
	Token        string
	IsAdmin      bool
}
