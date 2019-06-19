package db

import (
	"context"
	"github.com/fresh-from-the-farm/authn/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const (
	DataBaseName          = "3FDataCore"
	AccountCollectionName = "Accounts"
	DataBaseURL           = "mongodb://localhost:27017"
)

var client *mongo.Client

func init() {
	log.Printf("Trying to connect DB-URI %v DB-NAME %v", DataBaseURL, DataBaseName)
	var err error
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI(DataBaseURL)
	client, err = mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Panicf("Failed to connect DB with %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Panicf("Failed to ping DB with %v", err)
	} else {
		log.Printf("Connected to DB-URI %v DB-NAME %v", DataBaseURL, DataBaseName)
	}
}

func AddAccount(account model.Account) error {
	collection := client.Database(DataBaseName).Collection(AccountCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, account)

	if err != nil {
		return err
	}

	log.Printf("Account created with ID %v and result %v", account.Username, result)
	return nil
}

func GetAccount(username string) (a model.Account, err error) {

	filter := bson.D{{"userName", username}}

	collection := client.Database(DataBaseName).Collection(AccountCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = collection.FindOne(ctx, filter).Decode(&a)

	if err != nil {
		log.Printf("Account lookup for username %v %v failed with %v", username, a, err)
		return a, err
	}

	return a, err
}

func UpdateAccount(account, update model.Account) (updated model.Account, err error) {
	collection := client.Database(DataBaseName).Collection(AccountCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = collection.UpdateOne(ctx, account, update)
	return update, err
}
