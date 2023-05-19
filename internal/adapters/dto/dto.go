package db_dto

type CreateUserDTO struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Email        string `json:"email" bson:"email"`
	Username     string `json:"username" bson:"username"`
	PasswordHash string `json:"-" bson:"password"`
}

type IsUserExists struct {
	Email    string `json:"email" bson:"email"`
	Username string `json:"username" bson:"username"`
}

type AuthByEmail struct {
	Email        string `json:"email" bson:"email"`
	PasswordHash string `json:"-" bson:"password"`
}

type AuthByUsername struct {
	Username     string `json:"username" bson:"username"`
	PasswordHash string `json:"-" bson:"password"`
}
