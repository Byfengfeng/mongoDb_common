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
	Use string 				`json:"use"`
	PassWd string 			`json:"pass_wd"`
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
func NewMongoDbSession(option *ConnectOption) *mgo.Session {
	//info := &mgo.DialInfo{
	//	Addrs:    []string{"mongodb://localhost:27017"},Timeout:  60 * time.Second,Database: option.Db,Username: option.Use,Password: option.PassWd}
	//session,err := mgo.DialWithInfo(info)
	session,err := mgo.Dial(fmt.Sprintf("%s:%d",option.Address,option.Port))
	//session, err := mgo.Dial(fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s",
	//	option.Use,
	//	option.PassWd,
	//	option.Address,
	//	option.Port,
	//	option.Db,
	//	option.Db))
	if err != nil {
		panic(err)
	}
	return session
}

func NewMongoDb(option *ConnectOption) inter.MongoDbInterface {
	session := NewMongoDbSession(option)
	//credential := &mgo.Credential{
	//		Mechanism: SHA1,
	//		Source: option.Db,
	//		Username: option.Use,
	//		Password: option.PassWd,
	//		}

	//session.Login(credential)
	db := session.DB(option.Db)
	db.Login(option.Use,option.PassWd)
	LogDb := db.C("test")
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

	query := bson.M{
		"$and": []bson.M{
			bson.M{"cid": bson.M{"eq": cid}},
			bson.M{"create_time": bson.M{"gt": startTime}},
			bson.M{"create_time": bson.M{"lt": endTime}},
		},
	}

	m.LogDb.Find(query).All(&logList) // 某个时间段
	return logList,nil
}