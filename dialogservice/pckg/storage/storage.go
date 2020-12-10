package storage

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database int

const (
	DBPostgres Database = 1
	DBMySql    Database = 2
)

const (
	// Username, Password, Proto, Address, Port, DBName
	DSNMySQL string = "%s:%s@%s(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"

	// Username, Password, DBName, Port, SSLMode
	DSNPostgres string = "user=%s password=%s dbname=%s port=%s sslmode=%s"
)

type Storage struct {
	*gorm.DB
}

func CreateConnection(db Database, template string, config *gorm.Config, params ...string) *Storage {
	interfaces := make([]interface{}, len(params))
	for i, v := range params {
		interfaces[i] = v
	}
	dsn := fmt.Sprintf(template, interfaces...)
	if config == nil {
		config = &gorm.Config{}
	}
	var openConnection *gorm.DB
	switch db {
	case DBPostgres:
		connect, err := gorm.Open(postgres.Open(dsn), config)
		if err != nil {
			panic(err)
		}
		openConnection = connect
	case DBMySql:
		connect, err := gorm.Open(mysql.Open(dsn), config)
		if err != nil {
			panic(err)
		}
		openConnection = connect
	default:
		panic("Non valid DB type. ")
	}
	storage := new(Storage)
	storage.DB = openConnection
	return storage
}


func (storage *Storage) CloseConnection() error {
	sqldb, err := storage.DB.DB()
	if err != nil {
		return err
	}
	err = sqldb.Close()
	if err != nil {
		return err
	}
	return nil
}