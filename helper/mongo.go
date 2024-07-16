package helper

import (
	"context"
	"errors"
	"fmt"

	"github.com/Barokah-AI/BackEnd/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// func InsertOneDoc(db *mongo.Database, col string, doc any) (insertedID primitive.ObjectID, err error) {
// 	result, err := db.Collection(col).InsertOne(context.Background(), doc)
// 	if err != nil {
// 		return
// 	}
// 	return result.InsertedID.(primitive.ObjectID), nil
// }


