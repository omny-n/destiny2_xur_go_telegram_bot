package database

import (
	"context"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DB         = "bot_users"
	COLLECTION = "users"
)

var clientInstance *mongo.Client
var clientInstanceError error
var mongoOnce sync.Once

type User struct {
	UserID    int    `bson:"_id"`
	FirstName string `bson:"first_name,omitempty"`
	LastName  string `bson:"last_name,omitempty"`
	UserName  string `bson:"user_name,omitempty"`
}

func Connect() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(os.Getenv("DATABASE"))
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		}
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		clientInstance = client
	})
	return clientInstance, clientInstanceError
}

func Insert(u User) error {
	client, err := Connect()
	if err != nil {
		log.Println(err)
		return err
	}
	collection := client.Database(DB).Collection(COLLECTION)
	_, err = collection.InsertOne(context.TODO(), u)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetUserById(id int) (User, error) {
	result := User{}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	client, err := Connect()
	if err != nil {
		return result, err
	}
	collection := client.Database(DB).Collection(COLLECTION)
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func GetAllUsers() ([]User, error) {
	filter := bson.D{{}}
	users := []User{}
	client, err := Connect()
	if err != nil {
		return users, err
	}
	collection := client.Database(DB).Collection(COLLECTION)
	cur, findError := collection.Find(context.TODO(), filter)
	if findError != nil {
		return users, findError
	}
	for cur.Next(context.TODO()) {
		t := User{}
		err := cur.Decode(&t)
		if err != nil {
			return users, err
		}
		users = append(users, t)
	}
	cur.Close(context.TODO())
	if len(users) == 0 {
		return users, mongo.ErrNoDocuments
	}
	return users, nil
}
