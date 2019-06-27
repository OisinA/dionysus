package services

import (
	"dionysus/models"

	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SubmissionService struct {
	Service
}

func (*SubmissionService) Add(submission_id string, user_id string, problem_id string, status int) error {
	elem := struct {
		Submission_ID string    `bson:"Submission_ID"`
		User_ID       string    `bson:"User_ID"`
		Problem_ID    string    `bson:"Problem_ID"`
		Status        int       `bson:"Status"`
		Updated       time.Time `bson:"Updated"`
	}{submission_id, user_id, problem_id, status, time.Now()}
	_, err := client.Database("dionysus").Collection("submissions").InsertOne(context.TODO(), elem)
	return err
}

func (*SubmissionService) List(params SearchParams) ([]*models.Submission, error) {
	collection := client.Database("dionysus").Collection("submissions")
	cur, err := collection.Find(context.TODO(), bson.M(params.Queries), &options.FindOptions{Sort: bson.D{{"Updated", -1}}})
	if err != nil {
		return nil, err
	}
	submissions := make([]*models.Submission, 0)
	for cur.Next(context.TODO()) {
		elem := struct {
			Submission_ID string    `bson:"Submission_ID"`
			User_ID       string    `bson:"User_ID"`
			Problem_ID    string    `bson:"Problem_ID"`
			Status        int       `bson:"Status"`
			Updated       time.Time `bson:"Updated"`
		}{}
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, &models.Submission{elem.Submission_ID, elem.User_ID, elem.Problem_ID, elem.Status, elem.Updated.Format("Mon Jan _2 15:04:05 2006")})
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	cur.Close(context.TODO())

	return submissions, nil
}
