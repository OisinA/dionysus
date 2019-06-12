package services

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService struct {
	Service
}

func (*AuthService) GetPasswordHash(username string) (string, error) {
	collection := client.Database("dionysus").Collection("users")
	elem := struct {
			ID       *primitive.ObjectID `bson:"_id,omitempty"`
			Username string              `bson:"Username,omitempty"`
			Password string              `bson:"Password,omitempty"`
			Email    string              `bson:"Email,omitempty"`
			Team     int                 `bson:"Team,omitempty"`
	}{}
	err := collection.FindOne(context.TODO(), bson.D{{"Username", username}}).Decode(&elem)
	if err != nil {
		return "", err
	}
	return elem.Password, nil
}