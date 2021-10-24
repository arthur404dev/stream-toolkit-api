package restream

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func storeTokens(t *TokenResponse) (string, error) {
	logger := log.WithFields(log.Fields{"source": "restream.storeTokens()", "access-token": t.AccessToken, "refresh-token": t.RefreshToken})
	logger.Debugln("token store started")
	if t.AccessToken == "" {
		logger.Errorln("no data was received from TokenResponse")
		return "", errors.New("no data was received inside of TokenResponse")
	}
	clientOptions := options.Client().
		ApplyURI(os.Getenv("MONGO_CREDENTIALS"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logger.Debugln("started db connection")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Errorln(err)
		return "", err
	}
	defer client.Disconnect(ctx)
	logger.Debugln("getting token collection and storage document")
	tokensCollection := client.Database("auth").Collection("tokens")
	id, err := primitive.ObjectIDFromHex(os.Getenv("MONGO_TOKEN_ID"))
	if err != nil {
		logger.Errorln(err)
	}
	logger.Debugln("replacing tokens")
	result, err := tokensCollection.ReplaceOne(ctx, bson.M{"_id": id}, *t)
	if err != nil {
		logger.Errorln(err)
		return "", err
	}

	logger.WithFields(log.Fields{"modified": result.ModifiedCount, "id": id}).Infoln("tokens replaced successfully")
	response := fmt.Sprintf("Tokens successfully updated on database at %+v", time.Now().Format(time.RFC3339))
	logger.Debugln("token store finished")
	return response, nil
}
func getTokens() (TokenResponse, error) {
	logger := log.WithFields(log.Fields{"source": "restream.getTokens()"})
	logger.Debugln("token fetch started")
	tr := TokenResponse{}
	clientOptions := options.Client().
		ApplyURI(os.Getenv("MONGO_CREDENTIALS"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logger.Debugln("started db connection")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Errorln(err)
		return tr, err
	}
	defer client.Disconnect(ctx)
	logger.Debugln("getting token collection and storage document")
	tokensCollection := client.Database("auth").Collection("tokens")
	id, err := primitive.ObjectIDFromHex(os.Getenv("MONGO_TOKEN_ID"))
	if err != nil {
		logger.Errorln(err)
	}
	logger.Debugln("fetching tokens")
	err = tokensCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&tr)
	if err != nil {
		logger.Errorln(err)
		return tr, err
	}
	logger.WithFields(log.Fields{"access-token": tr.AccessToken, "refresh-token": tr.RefreshToken}).Infoln("tokens retrieved successfully")
	logger.Debugln("token fetch finished")
	return tr, nil
}
