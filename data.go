package data

type Log struct {
	Cid int64 `json:"cid"`
	LogLv int32 `json:"log_lv"`
	Content string `json:"content"`
	CreateTime int64 `json:"create_time"`
}

type LogEnum string

const (
	Cid = "cid"
	LogLv = "log_lv"
	Content ="context"
	CreateTime = "create_time"
)