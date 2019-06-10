package controllers

import (
	"github.com/bwmarrin/lit"
)

type ChallengeController struct{}

func (*ChallengeController) ValidateID(id int) int {
	if id < 0 {
		lit.Error("ID incorrect")
	}
	return id
}
