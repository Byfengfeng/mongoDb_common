package common

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/mgo.v2"
	data "mongoDb_common"
	"mongoDb_common/inter"
)

//mongoDb config data
type ConnectOption struct {
	Address string  		`json:"address"`
	Port uint16 			`json:"port"`
	MongoDbUse string 		`json:"mongo_db_use"`
	MongoDbPassWd string 	`json:"mongo_db_pass_wd"`
	DbUse string 			`json:"use"`
	DbPassWd string 		`json:"pass_wd"`
	Db string 				`json:"db"`
	MaxConnPoolSize uint64  `json:"max_conn_pool_size"`
}

type Encryption string

const (
	SHA1 = "SCRAM-SHA-1"
	SHA256 = "SCRAM-SHA-256"
)

type MongoDb struct {
	Session *mgo.Session
	LogDb *mgo.Collection
}



//url: mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb
func NewMongoDb(option *ConnectOption) inter.MongoDbInterface {
	session,err := mgo.Dial(fmt.Sprintf("%s:%s@%s:%d",option.MongoDbUse,option.MongoDbPassWd,option.Address,option.Port))
	if err != nil {
		panic(err)
	}

	//session.SetMode(mgo.Monotonic, true)
	db := session.DB(option.Db)
	err = db.Login(option.DbUse,option.DbPassWd)
	if err != nil {
		panic(err)
	}
	session.SetPoolLimit(100)
	LogDb := db.C(option.Db)
	defer session.Close()
	return &MongoDb{Session: session,LogDb: LogDb}
}

func (m *MongoDb) GetLogCollection(dataBaseName,tableName string) *mgo.Collection {
	m.LogDb = m.Session.DB(dataBaseName).C(tableName)
	return m.LogDb
}

func (m *MongoDb) GetDatabase(dataBaseName string) *mgo.Database {
	return m.Session.DB(dataBaseName)
}

func (m *MongoDb) GetCollection(dataBaseName,tableName string) *mgo.Collection {
	return m.Session.DB(dataBaseName).C(tableName)
}

func (m *MongoDb) AddLog(log interface{})  {
	err := m.LogDb.Insert(&log)
	if err != nil {
		panic(err)
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

func (m *MongoDb) FindLog(cid,startTime,endTime int64) ([]data.Log,error) {
	iter := m.LogDb.Find(bson.D{
		//{"cid", cid},
		//{"$gte", bson.A{"create_time", startTime}},
		//{"$lte", bson.A{"create_time", endTime}},
	}).Sort("createtime").Iter()
	logs := make(map[string]interface{},0)
	logData := make([]data.Log,0)
	for iter.Next(&logs) {
		if len(logs) > 0 {
			logData = append(logData,data.Log{
				Cid: logs[data.Cid].(int64),
			})
		}

	}
	if err := iter.Close(); err != nil {
		return nil,err
	}
	return logData,nil
}

func (m *MongoDb) FindOne(cid,startTime,endTime int64) ([]data.Log,error) {
	//log := data.Log{}
	logList := make([]data.Log,0)

	//query := bson.M{
	//	//bson.M{
	//	"cid": bson.M{"$eq": cid},
	//		"create_time": bson.M{"$gt": startTime,"$lt": endTime},
	//		//bson.M{"$lte": endTime},
	//	//},
	//}
	var query []bson.M
	query = append(query,bson.M{
		"cid": bson.M{"$eq": cid},
	})
	query = append(query,bson.M{
		"createtime": bson.M{"$gte": startTime},
	})
	query = append(query,bson.M{
		"createtime": bson.M{"$lt": endTime+1},
	})
	m.LogDb.Find(bson.M{"$and": query}).All(&logList) // 某个时间段
	return logList,nil
}