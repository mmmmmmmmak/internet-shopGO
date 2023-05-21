package entity

type User struct {
	ID           string  `json:"id" bson:"_id,omitempty"`
	Email        string  `json:"email" bson:"email"`
	Username     string  `json:"username" bson:"username"`
	PasswordHash string  `json:"-" bson:"password"`
	Session      Session `json:"session" bson:"session"`
}

type Session struct {
	RefreshToken string `json:"refreshToken" bson:"refreshToken"`
	ExpiresAt    int64  `json:"expiresAt" bson:"expiresAt"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken" bson:"accessToken"`
	RefreshToken string `json:"refreshToken" bson:"refreshToken"`
}
