package dto

type CreateProductDTO struct {
	Name        string  `json:"name" bson:"name"`
	Description string  `json:"description" bson:"description"`
	Price       float64 `json:"price" bson:"price"`
	Token       string  `json:"token" bson:"token"`
}
