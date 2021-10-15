package inter

import (
	"gopkg.in/mgo.v2"
	data "mongoDb_common"
)

type MongoDbInterface interface {
	GetDatabase(dataBaseName string) *mgo.Database
	GetCollection(dataBaseName,tableName string) *mgo.Collection

	AddLog(log interface{})
	FindLog(cid,startTime,endTime int64) ([]data.Log,error)
	FindOne(cid,startTime,endTime int64) ([]data.Log,error)
}


