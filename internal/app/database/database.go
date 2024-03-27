package database

import (
	"context"

	"github.com/pheezz/medods-jwt-test/internal/app/config"

	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log"
)

var conf = config.Conf

var errUserNotFound = errors.New("user not found")

type KeySchema struct {
	RefreshTokenHash  []byte `bson:"refresh_token_hash"`
	ExpireAtTimestamp int64  `bson:"expire_at_timestamp"`
	FromIP            string `bson:"from_ip"`
	Fingerprint       string `bson:"fingerprint"`
}

type BaseUserSchema struct {
	GUID string      `bson:"_id"`
	Keys []KeySchema `bson:"keys"`
}

func GetUserByGUID(GUID string) (BaseUserSchema, error) {
	collection := getUserCollection()
	var result BaseUserSchema
	filter := bson.D{{Key: "_id", Value: GUID}}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return BaseUserSchema{}, errUserNotFound
		default:
			log.Panic(err)
		}
	}
	defer closeConnection(collection)

	return result, nil
}

func UpsertUser(user BaseUserSchema) {
	collection := getUserCollection()
	filter := bson.D{{Key: "_id", Value: user.GUID}}
	update := bson.D{{Key: "$set", Value: user}}
	_, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		log.Panic(err)
	}
	defer closeConnection(collection)
}

func AddKeyToUser(GUID string, key KeySchema) {
	collection := getUserCollection()
	filter := bson.D{{Key: "_id", Value: GUID}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "keys", Value: key}}}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Panic(err)
	}
	defer closeConnection(collection)
}

func GetUserByFingerprintAndIP(fingerprint string, ip string) (BaseUserSchema, KeySchema, error) {
	collection := getUserCollection()
	var user BaseUserSchema
	filter := bson.D{
		{Key: "keys.fingerprint", Value: fingerprint},
		{Key: "keys.from_ip", Value: ip},
	}
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return BaseUserSchema{}, KeySchema{}, errUserNotFound
		default:
			log.Panic(err)
		}
	}
	defer closeConnection(collection)

	var key KeySchema
	for _, k := range user.Keys {
		if k.Fingerprint == fingerprint && k.FromIP == ip {
			key = k
			break
		}
	}
	return user, key, nil

}

func DeleteKey(key KeySchema) {
	collection := getUserCollection()
	filter := bson.D{
		{Key: "keys.from_ip", Value: key.FromIP},
		{Key: "keys.refresh_token_hash", Value: key.RefreshTokenHash},
	}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "keys", Value: key}}}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Panic(err)
	}
	defer closeConnection(collection)
}

func getUserCollection() *mongo.Collection {
	return getDatabase().Collection("users")
}

func getDatabase() *mongo.Database {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(conf.MongoURI))
	if err != nil {
		log.Panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Panic(err)
	}

	return client.Database(conf.DatabaseName)

}

func closeConnection(collection *mongo.Collection) {
	err := collection.Database().Client().Disconnect(context.Background())
	if err != nil {
		log.Panic(err)
	}
}
