package config

import "github.com/zeromicro/go-zero/zrpc"

type MySQLConf struct {
	Host      string
	Port      int
	Database  string
	Username  string
	Password  string
	Charset   string
	Collation string
	Timeout   string
}

type Config struct {
	zrpc.RpcServerConf
	MySQLConf MySQLConf
}
