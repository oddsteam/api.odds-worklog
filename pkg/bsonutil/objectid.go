package bsonutil

import "go.mongodb.org/mongo-driver/bson/primitive"

func MustObjectIDFromHex(s string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return id
}

func ObjectIDFromHex(s string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(s)
}
