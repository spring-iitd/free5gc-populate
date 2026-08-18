package runtime

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// This function come from https://github.com/free5gc/MongoDBLibrary/blob/main/api_mongoDB.go (License Apache 2)
// with new parameters "client" and "dbName", and a change on the return type
func RestfulAPIPost(client *mongo.Client, dbName string, collName string, filter bson.M, postData map[string]interface{}) (bool, error) {
	if client == nil {
		return false, fmt.Errorf("mongo client is nil")
	}
	collection := client.Database(dbName).Collection(collName)

	var checkItem map[string]interface{}
	err := collection.FindOne(context.TODO(), filter).Decode(&checkItem)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return false, fmt.Errorf("find in %s: %w", collName, err)
	}

	if checkItem == nil {
		if _, err := collection.InsertOne(context.TODO(), postData); err != nil {
			return false, fmt.Errorf("insert into %s: %w", collName, err)
		}
		return false, nil
	} else {
		if _, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set": postData}); err != nil {
			return true, fmt.Errorf("update %s: %w", collName, err)
		}
		return true, nil
	}
}

// This function come from https://github.com/free5gc/MongoDBLibrary/blob/main/api_mongoDB.go (License Apache 2)
// with new parameters "client" and "dbName", and a change on the return type
func RestfulAPIPostMany(client *mongo.Client, dbName string, collName string, filter bson.M, postDataArray []interface{}) (bool, error) {
	if client == nil {
		return false, fmt.Errorf("mongo client is nil")
	}
	collection := client.Database(dbName).Collection(collName)

	if _, err := collection.InsertMany(context.TODO(), postDataArray); err != nil {
		return false, fmt.Errorf("insert many into %s: %w", collName, err)
	}
	return false, nil
}
