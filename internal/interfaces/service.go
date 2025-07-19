package interfaces

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=UserService --output=../../mocks
type UserService interface {
	RegisterUser(login, password string) error
}