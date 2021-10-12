package common

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
