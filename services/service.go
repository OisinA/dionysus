package services

import (
	"context"

	"github.com/bwmarrin/lit"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct{}

var client *mongo.Client

func Setup() error {
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
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
