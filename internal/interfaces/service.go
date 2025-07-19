package interfaces

type UserService interface {
	RegisterUser(login, password string) error
}