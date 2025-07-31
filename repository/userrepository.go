package repository

import (
	"context"

	"github.com/sol-tad/Blog-post-Api/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryImpl struct {
	collection *mongo.Collection
}
func NewUserRepository(coll *mongo.Collection) domain.UserRepository {
	return &UserRepositoryImpl{
		collection: coll,
	}
}


func (ur *UserRepositoryImpl) CreateUser(user domain.User) (domain.User,error){

		user.Role="user"
		_,err:=ur.collection.InsertOne(context.Background(),user)
		return user,err
}
func (ur *UserRepositoryImpl) Login(username string,password string)(domain.User,error){
	var user domain.User
    err := ur.collection.FindOne(context.Background(), map[string]string{"username": username, "password": password}).Decode(&user)
    return user, err
}