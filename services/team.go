package services

import (
	"dionysus/models"
	
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeamService struct {
	Service
}

func (*TeamService) List(params SearchParams) ([]*models.Team, error) {
	collection := client.Database("dionysus").Collection("teams")
	cur, err := collection.Find(context.TODO(), bson.M(params.Queries))
	if err != nil {
		return nil, err
	}
	teams := make([]*models.Team, 0)
	for cur.Next(context.TODO()) {
		elem := struct {
			ID        *primitive.ObjectID `bson:"_id,omitempty"`
			Team_Name string              `bson:"Team_Name,omitempty"`
		}{}
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		teams = append(teams, &models.Team{elem.ID.Hex(), elem.Team_Name})
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	cur.Close(context.TODO())

	return teams, nil
}

func (*TeamService) ListMembers(params SearchParams) ([]string, error) {
	collection := client.Database("dionysus").Collection("team_members")
	cur, err := collection.Find(context.TODO(), bson.M(params.Queries))
	if err != nil {
		return nil, err
	}
	users := make([]string, 0)
	for cur.Next(context.TODO()) {
		elem := struct {
			Team_ID string `bson:"Team_ID"`
			User_ID string `bson:"User_ID"`
		}{}
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		users = append(users, elem.User_ID)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	cur.Close(context.TODO())

	return users, nil
}