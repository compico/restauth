package user

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type Users struct {
	User []User
}

type User struct {
	GUID    string
	Access  []*jwt.Token
	Refresh []*jwt.Token
	Tokens  []*mongo.InsertOneResult
}

func InitUsers() *Users {
	return new(Users)
}

func (users Users) NewUser(guid string) {
	user := User{
		GUID: guid,
	}
	users.User = append(users.User, user)
}
