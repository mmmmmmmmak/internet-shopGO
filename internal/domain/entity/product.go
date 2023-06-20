package entity

type Product struct {
	ID          string  `json:"id" bson:"_id,omitempty"`
	Name        string  `json:"name" bson:"name"`
	Description string  `json:"description" bson:"description"`
	Price       float64 `json:"price" bson:"price"`
	Discount    float64 `json:"discount" bson:"discount"`
	Seller      string  `json:"seller" bson:"seller"`
	Rating      float64 `json:"rating" bson:"rating"`
}
