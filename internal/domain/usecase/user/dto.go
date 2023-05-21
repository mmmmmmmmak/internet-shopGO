package user_usecase

type CreateUserDTO struct {
	Email    string `json:"email" bson:"email"`
	Username string `json:"username" bson:"username"`
	Password string `json:"-" bson:"password"`
}

type AuthUserDTO struct {
	Email    string `json:"email" bson:"email"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken" bson:"accessToken"`
	RefreshToken string `json:"refreshToken" bson:"refreshToken"`
}

type GetUserDTO struct {
	Token string `json:"token" bson:"token"`
}

type RefreshTokenDTO struct {
	Token string `json:"refreshToken" bson:"refreshToken"`
}
