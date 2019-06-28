package services

import (
	"context"
	"os"

	"github.com/bwmarrin/lit"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct{}

type SearchParams struct {
	Queries map[string]interface{}
}

var client *mongo.Client

func Setup() error {
	clientOptions := options.Client().ApplyURI("mongodb://" + os.Getenv("DIONYSUS_DB") + ":27017")
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	lit.Info("Connected to database.")
	return nil
}

func Cleanup() {
	client.Disconnect(context.TODO())
	lit.Info("Disconnected from database.")
}
