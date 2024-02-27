package localfile

import (
	"fmt"
	"github.com/skruger/privatestudio/web/datastore"
)

type userDatabase struct {

}

func NewLocalfileUserDatabase() datastore.Login {
	return userDatabase{}
}

type user struct {
	username string
	groups []string
}

func (u userDatabase) CheckLogin(username string, password string) (datastore.UserInfo, error) {
	return &user{
		username: username,
	}, nil
}

func (u userDatabase) CreateUser(username string) (datastore.UserInfo, error) {
	return &user{
		username: username,
	}, nil
}

func (u userDatabase) GetUser(username string) (datastore.UserInfo, error) {
	return &user{
		username: username,
	}, nil
}

func (u user) GetUsername() string {
	return u.username
}

func (u user) SetPassword(password string) error {
	return fmt.Errorf("SetPasswrod not implemented for localfile users")
}

func (u user) GetGroups() []string {
	return u.groups
}
