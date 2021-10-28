package inter

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDbInterface interface {
	GetDatabase(dataBaseName string) *mongo.Database
	GetCollection(dataBaseName,tableName string) *mongo.Collection

	AddLog(cid,createTime int64,logLv int8,log string)
	FindLog(cid,startTime,endTime int64) ([]map[string]interface{},error)
}


