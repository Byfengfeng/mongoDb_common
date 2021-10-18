package main

import (
	"fmt"
	data "mongoDb_common"
	"mongoDb_common/common"
	"time"
)

func main() {
	options := &common.ConnectOption{
		Address: "localhost",
		Port: 27017,
		Use: "aonuo",
		PassWd: "8888",
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