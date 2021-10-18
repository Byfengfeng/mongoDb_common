package main

import (
	"fmt"
	data "mongoDb_common"
	"mongoDb_common/common"
	"time"
)

type Log struct {
	Cid int64
	CreateTime int64
	Log string
}

func main() {
	options := &common.ConnectOption{
		Address: "127.0.0.1",
		Port: 27017,
		MongoDbUse: "admin",
		MongoDbPassWd: "123456",
		DbUse: "aonuo",
		DbPassWd: "8888",
		Db: "test",
		MaxConnPoolSize: 5,
	}
	dbSession := common.NewMongoDb(options)
	currentTime := time.Now().UnixNano() / 1e6
	dbSession.AddLog(data.Log{
		Cid: 1,
		LogLv: 2,
		Content: "123",
		CreateTime: currentTime,
	})
	//log, err := dbSession.FindLog(1, currentTime, currentTime)
	//if err != nil {
	//	panic(err)
	//}
	one, err := dbSession.FindOne(1, 0, currentTime)
	if err != nil {
		panic(err)
	}
	fmt.Println(one)
}