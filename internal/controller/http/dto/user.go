package dto

type CreateUserDTO struct {
	Email    string `json:"email" bson:"email"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type AuthDTO struct {
	Email    string `json:"email" bson:"email"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type GetUserDTO struct {
	Token string `json:"token" bson:"token"`
}

type RefreshTokenDTO struct {
	Token string `json:"refreshToken" bson:"refreshToken"`
}
