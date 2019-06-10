package services

import (
	"dionysus/models"

	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChallengeService struct {
	Service
}

func (*ChallengeService) List() ([]*models.Challenge, error) {
	collection := client.Database("dionysus").Collection("challenges")
	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	challenges := make([]*models.Challenge, 0)
	for cur.Next(context.TODO()) {
		elem := struct {
			ID           *primitive.ObjectID `json:"ID" bson:"_id,omitempty"`
			Challenge_ID int                 `bson:"ID,omitempty"`
			Name         string              `bson:"Name,omitempty"`
		}{}
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		challenges = append(challenges, &models.Challenge{ID: elem.Challenge_ID, Name: elem.Name})
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	cur.Close(context.TODO())

	return challenges, nil
}

func (*ChallengeService) Add(challenge models.Challenge) error {
	elem := struct {
		ID   int    `bson:"ID"`
		Name string `bson:"Name"`
	}{ID: challenge.ID, Name: challenge.Name}
	_, err := client.Database("dionysus").Collection("challenges").InsertOne(context.TODO(), elem)
	return err
}
