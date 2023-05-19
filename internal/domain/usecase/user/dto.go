package user_usecase

type CreateUserDTO struct {
	Email    string `json:"email" bson:"email"`
	Username string `json:"username" bson:"username"`
	Password string `json:"-" bson:"password"`
}

type AuthByEmail struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type AuthByUsername struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type AuthUser struct {
	Email    string `json:"email" bson:"email"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
