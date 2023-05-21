package db_dto

type CreateUserDTO struct {
	ID           string
	Email        string
	Username     string
	PasswordHash string
	Session      Session
}

type Session struct {
	RefreshToken string
	ExpiresAt    int64
}

type IsUserExistsDTO struct {
	Email    string
	Username string
}

type AuthUserDTO struct {
	Email        string
	Username     string
	PasswordHash string
}

type GetUserDTO struct {
	ID string
}

type GetUserByRefreshTokenDTO struct {
	Token string
}
