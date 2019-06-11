package controllers

import (
	"errors"
)

type ChallengeController struct{}

func (*ChallengeController) ValidateID(id int) (int, error) {
	if id < 0 {
		return -1, errors.New("ID not in valid range")
	}
	return id, nil
}
