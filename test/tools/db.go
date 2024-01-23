package toolstester

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// dbConfig dbConfig
type dbConfig struct {
	Type            string `yaml:"type"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Database        string `yaml:"database"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	TablePrefix     string `yaml:"table_prefix"`
	PoolMaxIdle     int    `yaml:"pool_max_idle"`
	PoolMaxActive   int    `yaml:"pool_max_active"`
	PoolIdleTimeout int    `yaml:"pool_idle_timeout"`
	ConnTimeout     int    `yaml:"conn_timeout"`
	ReadTimeout     int    `yaml:"read_timeout"`
	WriteTimeout    int    `yaml:"write_timeout"`
}

// GetDBConnection Get DB Connection
func GetDBConnection(c *dbConfig) *gorm.DB {
	if db != nil {
		return db
	}

	var (
		dbType, dbName, user, password, host, tablePrefix string
		port                                              int
	)
	dbType = c.Type
	dbName = c.Database
	user = c.User
	password = c.Password
	host = c.Host
	port = c.Port
	tablePrefix = c.TablePrefix
	db, err := gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		user,
		password,
		host,
		port,
		dbName)+"&loc=Asia%2fShanghai")

	if err != nil {
		log.Panic(err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	// db.SingularTable(true)
	// db.LogMode(true)
	db.DB().SetMaxIdleConns(c.PoolMaxIdle)
	db.DB().SetMaxOpenConns(c.PoolMaxActive)
	db.DB().SetConnMaxLifetime(time.Second * time.Duration(c.PoolIdleTimeout))

	return db
}

// Close Close db connection
func Close(db *gorm.DB) {
	db.Close()
}
