package storage

type Storage interface {
	GetUserByEmail(userEmail string) (string, error)
	NewUserRegister(userEmail, userName, userPasswordHash string) (int64, error)
}
