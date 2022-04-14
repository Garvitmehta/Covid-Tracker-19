package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type State struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	State     string             `json:"state,omitempty"`
	Confirmed float64            `json:"Confirmed,omitempty" validate:"required"`
	Deceased  float64            `json:"deceased,omitempty" validate:"required"`
}
