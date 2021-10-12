package mongoDb_common

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mongoDb_common/common"
	"time"
)

func NewMongoDbClient(option *common.ConnectOption) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%d",option.Address,option.Port)),
		GetOptions(option.Db,option.Use,option.PassWd,option.MaxConnPoolSize))
	if err != nil {
		panic(err)
	}
	return client
}

func GetOptions(db,user,passwd string,maxPoolSize uint64) *options.ClientOptions {
	opts:= &options.ClientOptions{
		Auth: &options.Credential{AuthMechanism:"SCRAM-SHA-1",
			AuthSource: db,
			Username: user,
			Password: passwd},
		MaxPoolSize: &maxPoolSize,
	}
	return opts
}

func GetDatabase(client *mongo.Client,dataBaseName string) *mongo.Database {
	return client.Database(dataBaseName)
}

func GetCollection(client *mongo.Client,dataBaseName,tableName string) *mongo.Collection {
	return client.Database(dataBaseName).Collection(tableName)
}