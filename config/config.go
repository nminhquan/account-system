package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	CFG_ROCKSDB_DIR          = "ROCKSDB_DIR"
	CFG_RAFT_LOG_DIR         = "RAFT_LOG_DIR"
	CFG_TC_SERVICE_HOST      = "TC_SERVICE_HOST"
	CFG_LOCK_SERVICE_HOST    = "LOCK_SERVICE_HOST"
	CFG_TC_NODE_SERVICE_HOST = "TC_NODE_SERVICE_HOST"
	CFG_RM_NODE_SERVICE_HOST = "RM_NODE_SERVICE_HOST"
	CFG_REDIS_HOST           = "REDIS_HOST"
)

/*
Stores config data
*/
var (
	RocksDBDir     string
	RaftLogDir     string
	TCServHost     string
	LockServHost   string
	TCNodeServHost string
	RMNodeServHost string
	RedisHost      string
)

func init() {
	fmt.Println("Init global config")
	LoadConfigFile("config", ".")
	fmt.Println(ToString())
}

func LoadConfigFile(fileName string, path string) {
	viper.SetConfigName(fileName)         // name of config file (without extension)
	viper.AddConfigPath(path)             // path to look for the config file in
	viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory
	err := viper.ReadInConfig()           // Find and read the config file
	if err != nil {                       // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	RocksDBDir = fmt.Sprintf("%v", viper.Get(CFG_ROCKSDB_DIR))
	RaftLogDir = fmt.Sprintf("%v", viper.Get(CFG_RAFT_LOG_DIR))
	TCServHost = fmt.Sprintf("%v", viper.Get(CFG_TC_SERVICE_HOST))
	LockServHost = fmt.Sprintf("%v", viper.Get(CFG_LOCK_SERVICE_HOST))
	TCNodeServHost = fmt.Sprintf("%v", viper.Get(CFG_TC_NODE_SERVICE_HOST))
	RMNodeServHost = fmt.Sprintf("%v", viper.Get(CFG_RM_NODE_SERVICE_HOST))
	RedisHost = fmt.Sprintf("%v", viper.Get(CFG_REDIS_HOST))
}

func ToString() string {
	return fmt.Sprintf("RocksDBDir = %v\nRaftLogDir = %v\nTCServHost = %v\nLockServHost = %v\nTCNodeServHost = %v\nRMNodeServHost = %v\nRedisHost = %v",
		RocksDBDir, RaftLogDir, TCServHost, LockServHost, TCNodeServHost, RMNodeServHost, RedisHost)
}
