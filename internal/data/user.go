// data/user.go
package data

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`  // json:"-" = never expose password
	Role         string `json:"role"`
}