package user

type User struct {
	Id           string `json:"id" bson:"_id, omitempty"`
	Username     string `json:"username" bson:"username"`
	PasswordHash string `json:"-" bson:"password"`
	Email        string `json:"email" bson:"email"`
}

type CreateUserDTO struct {
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Email        string `json:"email"`
}
