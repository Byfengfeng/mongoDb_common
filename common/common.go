package common

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	data "mongoDb_common"
	"mongoDb_common/inter"
	"time"
)

//mongoDb config data
type ConnectOption struct {
	Address string  		`json:"address"`
	Port uint16 			`json:"port"`
	Use string 				`json:"use"`
	PassWd string 			`json:"pass_wd"`
	Db string 				`json:"db"`
	MaxConnPoolSize uint64  `json:"max_conn_pool_size"`
}

const (
	SHA1 = "SCRAM-SHA-1"
	SHA256 = "SCRAM-SHA-256"
)

type MongoDb struct {
	Client *mongo.Client
	LogDb *mongo.Collection
}

func NewMongoDbClient(option *ConnectOption) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%d",option.Address,option.Port)),
		&options.ClientOptions{
			Auth: &options.Credential{AuthMechanism: SHA1,
				AuthSource: option.Db,
				Username: option.Use,
				Password: option.PassWd},
			MaxPoolSize: &option.MaxConnPoolSize,
		})
	if err != nil {
		panic(err)
	}
	return client
}




func NewMongoDb(option *ConnectOption) inter.MongoDbInterface {
	client := NewMongoDbClient(option)
	mongoDb := &MongoDb{Client: client}
	mongoDb.GetLogCollection(option.Db,"log")
	return mongoDb
}

func (m *MongoDb) GetLogCollection(dataBaseName,tableName string) *mongo.Collection {
	m.LogDb = m.Client.Database(dataBaseName).Collection(tableName)
	return m.LogDb
}

func (m *MongoDb) GetDatabase(dataBaseName string) *mongo.Database {
	return m.Client.Database(dataBaseName)
}

func (m *MongoDb) GetCollection(dataBaseName,tableName string) *mongo.Collection {
	return m.Client.Database(dataBaseName).Collection(tableName)
}

func (m *MongoDb) AddLog(cid,createTime int64,logLv int8,log string)  {
	_, err := m.LogDb.InsertOne(context.TODO(), bson.D{
		{"cid", cid},
		{"log_lv", logLv},
		{"context", log},
		{"create_time", createTime},
	})
	if err != nil {
		return
	}
}

/**
小于		{<key>:{$lt:<value>}}	db.col.find({"likes":{$lt:50}}).pretty()	where likes < 50
小于或等于	{<key>:{$lte:<value>}}	db.col.find({"likes":{$lte:50}}).pretty()	where likes <= 50
大于		{<key>:{$gt:<value>}}	db.col.find({"likes":{$gt:50}}).pretty()	where likes > 50
大于或等于	{<key>:{$gte:<value>}}	db.col.find({"likes":{$gte:50}}).pretty()	where likes >= 50
不等于		{<key>:{$ne:<value>}}	db.col.find({"likes":{$ne:50}}).pretty()	where likes != 50
 */

func (m *MongoDb) FindLog(cid,startTime,endTime int64) ([]*data.Log,error) {
	//cursor,err := m.LogDb.Find(context.TODO(),bson.D{
	//	{"cid", cid},
		//{"$gte", bson.A{"create_time", startTime}},
		//{"$lte", bson.A{"create_time", endTime}},
	//})
	cursor,err := m.LogDb.Find(context.TODO(),bson.M{
		"cid": cid,
		//{"$gte", bson.A{"create_time", startTime}},
		//{"$lte", bson.A{"create_time", endTime}},
	})
	if err != nil {
		return nil,err
	}
	logs := make([]*data.Log,0)
	for cursor.Next(context.TODO()) {
		logMap := make(map[string]interface{})
		err = cursor.Decode(&logMap)
		if err != nil {
			panic(err)
		}
		if len(logMap) > 0 {
			logs = append(logs,&data.Log{
				Cid: logMap[data.Cid].(int64),
				LogLv: logMap[data.LogLv].(int32),
				Content: logMap[data.Content].(string),
				CreateTime: logMap[data.CreateTime].(int64),
			})
		}

	}
	return logs,nil
}


