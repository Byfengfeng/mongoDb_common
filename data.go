package data

type Log struct {
	Id string `json:"_id"`
	Cid int64 `json:"cid"`
	LogLv int32 `json:"log_lv"`
	Content string `json:"content"`
	CreateTime int64 `json:"create_time"`
}

type LogEnum string

const (
	Id = "_id"
	Cid = "cid"
	LogLv = "log_lv"
	Content ="context"
	CreateTime = "create_time"
)