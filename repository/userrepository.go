package repository

import (
	"context"

	"github.com/sol-tad/Blog-post-Api/domian"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryImpl struct {
	collection *mongo.Collection
}
func NewUserRepository(coll *mongo.Collection) domian.UserRepository {
	return &UserRepositoryImpl{
		collection: coll,
	}
}


func (ur *UserRepositoryImpl) CreateUser(user domian.User) (domian.User,error){

		user.Role="user"
		_,err:=ur.collection.InsertOne(context.Background(),user)
		return user,err
}
func (ur *UserRepositoryImpl) Login(username string,password string)(domian.User,error){
	var user domian.User
    err := ur.collection.FindOne(context.Background(), map[string]string{"username": username, "password": password}).Decode(&user)
    return user, err
}