package main

import (
	"fmt"
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
	mongoDb := common.NewMongoDb(options)
	currentTime := time.Now().UnixNano() / 1e6
	mongoDb.AddLog(1,currentTime,2,"asdasdasd")
	logs, err := mongoDb.FindLog(1, currentTime, currentTime)
	if err != nil {
		panic(err)
	}
	if len(logs) > 0 {
		for _,log := range logs {
			fmt.Println(log)
		}
	}

}