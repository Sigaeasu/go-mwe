package config

import (
	"path/filepath"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Environment string = "dev"
var Config config

type config struct {
	Environment string `mapstructure:"environment"`
	PostgresCfg struct {
		Database    string `mapstructure:"database"`
		Host        string `mapstructure:"host"`
		Port        string `mapstructure:"port"`
		Username    string `mapstructure:"username"`
		Password    string `mapstructure:"password"`
		MaxConn     int    `mapstructure:"max_conn"`
		MinIdleConn int    `mapstructure:"min_idle_conn"`
		MaxRetries  int    `mapstructure:"max_retries"`
	} `mapstructure:"postgres"`
	JWTCfg struct {
		Issuer  string `mapstructure:"issuer"`
		Exp     int    `mapstructure:"exp"`
		SignKey string `mapstructure:"sign_key"`
	} `mapstructure:"jwt"`
}

func init() {
	var err error

	viper.SetConfigName("application."+Environment)
	viper.AddConfigPath(filepath.Join(GetAppBasePath(), "config/app"))
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
			logrus.Errorf("Error on config file: %v", err)
		}
	err = viper.Unmarshal(&Config)
	if err != nil {
		logrus.Errorf("Fail on decoding to struct, %v", err)
	}
}

func GetAppBasePath() string {
	basePath, _ := filepath.Abs(".")
	for filepath.Base(basePath) != "go-mwe" {
		basePath = filepath.Dir(basePath)
	}

	return basePath
}