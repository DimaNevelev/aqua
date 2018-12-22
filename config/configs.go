package config

import (
	"database/sql"
	"github.com/dimanevelev/aqua/persistence"
	_ "github.com/go-sql-driver/mysql"
)

type ServerConfig struct {
	Database *sql.DB
	ServerConst
}

type ServerConst struct {
	Port      string
	MySqlConf persistence.MySqlConf
}

// NewConfig is used to generate a configuration instance which will be passed around the codebase
func NewServerConfig(consts ServerConst) (ServerConfig, error) {
	var config ServerConfig
	config.ServerConst = consts
	var err error
	config.Database, err = persistence.InitClient(consts.MySqlConf)
	if err != nil {
		return config, err
	}
	return config, err
}

type TraverserConfig struct {
	Url     string
	Path    string
	Threads int
}
