package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/bwmarrin/lit"
)

type ScoreService struct {
	Service
}

func (*ScoreService) Add(submission_id string, user_id string, problem_id string, score int) error {
	elem := struct {
		Submission_ID string `bson:"Submission_ID"`
		User_ID    string `bson:"User_ID"`
		Problem_ID string `bson:"Problem_ID"`
		Score      int    `bson:"Score"`
	}{submission_id, user_id, problem_id, score}
	_, err := client.Database("dionysus").Collection("scores").InsertOne(context.TODO(), elem)
	return err
}

func (*ScoreService) Get(user_id string) (map[string]int, error) {
	collection := client.Database("dionysus").Collection("scores")
	cur, err := collection.Find(context.TODO(), bson.D{{"User_ID", user_id}})
	if err != nil {
		return nil, err
	}
	scores := make(map[string]int, 0)
	for cur.Next(context.TODO()) {
		elem := struct {
			User_ID    string `json:"user_id" bson:"User_ID"`
			Problem_ID string `json:"problem_id" bson:"Problem_ID"`
			Score      int    `json:"score" bson:"Score"`
		}{}
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		if val, ok := scores[elem.Problem_ID]; ok {
			if elem.Score > val {
				scores[elem.Problem_ID] = elem.Score
			}
			continue
		}
		scores[elem.Problem_ID] = elem.Score
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	cur.Close(context.TODO())

	return scores, nil
}

func (*ScoreService) GetSubmissionScore(submission_id string) (int, error) {
	collection := client.Database("dionysus").Collection("scores")
	elem := struct {
		Submission_ID string `bson:"Submission_ID"`
		User_ID       string `bson:"User_ID"`
		Problem_ID    string `bson:"Problem_ID"`
		Score         int    `bson:"Score"`
	}{}
	lit.Debug(submission_id)
	err := collection.FindOne(context.TODO(), bson.D{{"Submission_ID", submission_id}}).Decode(&elem)
	if err != nil {
		return -1, err
	}
	return elem.Score, nil
}