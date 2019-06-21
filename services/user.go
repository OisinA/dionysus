package services

import (
	"dionysus/models"
	
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/bwmarrin/lit"
	"fmt"
)

type UserService struct {
	Service
}

func (*UserService) List(params SearchParams) ([]*models.User, error) {
	collection := client.Database("dionysus").Collection("users")
	cur, err := collection.Find(context.TODO(), bson.M(params.Queries))
	if err != nil {
		return nil, err
	}
	users := make([]*models.User, 0)
	for cur.Next(context.TODO()) {
		elem := struct {
			ID       *primitive.ObjectID `bson:"_id,omitempty"`
			Username string              `bson:"Username,omitempty"`
			Password string              `bson:"Password,omitempty"`
			Email    string              `bson:"Email,omitempty"`
		}{}
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		users = append(users, &models.User{elem.ID.Hex(), elem.Username, elem.Password, elem.Email})
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	cur.Close(context.TODO())

	return users, nil
}

func (*UserService) Get(id string) (models.User, error) {
	collection := client.Database("dionysus").Collection("users")
	elem := struct {
			ID       *primitive.ObjectID `bson:"_id,omitempty"`
			Username string              `bson:"Username,omitempty"`
			Password string              `bson:"Password,omitempty"`
			Email    string              `bson:"Email,omitempty"`
	}{}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}
	err = collection.FindOne(context.TODO(), bson.D{{"_id", objectID}}).Decode(&elem)
	return models.User{elem.ID.Hex(), elem.Username, elem.Password, elem.Email}, err
}

func (*UserService) Add(user models.User) error {
	elem := struct {
			Username string              `bson:"Username,omitempty"`
			Password string              `bson:"Password,omitempty"`
			Email    string              `bson:"Email,omitempty"`
	}{user.Username, user.Password, user.Email}
	lit.Debug(fmt.Sprint(elem))
	_, err := client.Database("dionysus").Collection("users").InsertOne(context.TODO(), elem)
	return err
}

func (*UserService) UsernameToID(username string) (string, error) {
	collection := client.Database("dionysus").Collection("users")
	elem := struct {
			ID       *primitive.ObjectID `bson:"_id,omitempty"`
			Username string              `bson:"Username,omitempty"`
			Password string              `bson:"Password,omitempty"`
			Email    string              `bson:"Email,omitempty"`
	}{}
	err := collection.FindOne(context.TODO(), bson.D{{"Username", username}}).Decode(&elem)
	if err != nil {
		return "", err
	}
	return elem.ID.Hex(), nil
}