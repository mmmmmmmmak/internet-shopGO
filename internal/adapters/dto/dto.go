package db_dto

type CreateUserDTO struct {
	ID           string
	Email        string
	Username     string
	PasswordHash string
}

type IsUserExists struct {
	Email    string
	Username string
}

type AuthUser struct {
	Email        string
	Username     string
	PasswordHash string
}
