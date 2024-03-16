package models

type User struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// type Product struct {
// 	Id       uint64  `json:"id"`
// 	Name     string  `json:"name"`
// 	Price    float64 `json:"price"`
// 	Quantity int     `json:"quantity"`
// }

// type Product_dto struct {
// 	Id    uint64  `json:"id"`
// 	Name  string  `json:"name"`
// 	Price float64 `json:"price"`
// }

func ToUserResponse(users []User) []UserResponse {
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
		})
	}
	return userResponses
}

// func ToProductResponse(products []Product) []Product_dto {
// 	var productResponses []Product_dto
// 	for _, product := range products {
// 		productResponses = append(productResponses, Product_dto{
// 			Name:  product.Name,
// 			Price: product.Price,
// 		})
// 	}
// 	return productResponses
// }
