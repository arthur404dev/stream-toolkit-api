package restream

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func storeTokens(t *TokenResponse) (string, error) {
	if t.AccessToken == "" {
		log.Fatalf("storeTokens error=No data was received inside of TokenResponse")
		return "", errors.New("No Data was received inside of TokenResponse")
	}
	clientOptions := options.Client().
		ApplyURI(os.Getenv("MONGO_CREDENTIALS"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("storeTokens.Mongo.Connect() error=%+v", err)
		return "", err
	}
	defer client.Disconnect(ctx)

	tokensCollection := client.Database("auth").Collection("tokens")
	id, err := primitive.ObjectIDFromHex(os.Getenv("MONGO_TOKEN_ID"))
	if err != nil {
		log.Fatalf("storeTokens.primitive.ObjectIDFromHex() error=%+v\n", err)
	}

	result, err := tokensCollection.ReplaceOne(ctx, bson.M{"_id": id}, t)
	if err != nil {
		log.Fatalf("storeTokens.tokensCollection.ReplaceOne() error=%+v\n", err)
		return "", err
	}

	log.Printf("StoreTokens success, modified documents=%+v, id=%+v raw response=%+v\n", result.ModifiedCount, id, result)
	response := fmt.Sprintf("Tokens successfully updated on database at %+v", time.Now().Format(time.RFC3339))

	return response, nil
}
func getTokens() (TokenResponse, error) {
	tr := TokenResponse{}
	clientOptions := options.Client().
		ApplyURI(os.Getenv("MONGO_CREDENTIALS"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("getTokens.Mongo.Connect() error=%+v", err)
		return tr, err
	}
	defer client.Disconnect(ctx)

	tokensCollection := client.Database("auth").Collection("tokens")
	id, err := primitive.ObjectIDFromHex(os.Getenv("MONGO_TOKEN_ID"))
	if err != nil {
		log.Fatalf("getTokenws.primitive.ObjectIDFromHex() error=%+v\n", err)
	}

	err = tokensCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&tr)
	if err != nil {
		log.Fatalf("getTokens.tokensCollection.ReplaceOne() error=%+v\n", err)
		return tr, err
	}
	return tr, nil
}
