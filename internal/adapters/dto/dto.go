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

type AuthByEmail struct {
	Email        string
	PasswordHash string
}

type AuthByUsername struct {
	Username     string
	PasswordHash string
}

type AuthUser struct {
	Email        string
	Username     string
	PasswordHash string
}
