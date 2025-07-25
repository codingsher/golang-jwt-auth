package types

type UserLogin struct {
	Id       int64  `json:"id"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserRegister struct {
	Id           int64  `json:"id"`
	UserName     string `json:"username" validate:"required"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}
