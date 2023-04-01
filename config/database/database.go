package database

import (
	pg "github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
)

type ParametersConnection struct {
	Username string
	Password string
	Host string
	Port string
	Database string
	MaxConnection int
	MinIdleConnection int
	MaxRetries int
}

func DatabaseConnection(params ParametersConnection) *pg.DB {
	logrus.Info("Start Database Connection")

	db := pg.Connect(&pg.Options{
		User: params.Username,
		Password: params.Password,
		Addr: params.Host+":"+params.Port,
		Database: params.Database,
		PoolSize: params.MaxConnection,
		MinIdleConns: params.MinIdleConnection,
		MaxRetries: params.MaxRetries,
	})

	return db
}