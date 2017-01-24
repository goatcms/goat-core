package db_manager

/*
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goatcms/cms-core/app/config"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

const (
	DbName = "db"
)

func DbFactory(dp dependency.Provider) (interface{}, error) {
	ins := dp.Get(config.ConfigName)
	config := ins.(config.Config)
	dbConfig := config.Database
	switch strings.ToLower(dbConfig.Adapter) {
	case "mysql":
		return CreateMyslq(dbConfig)
	case "pgsql":
		return CreatePgsql(dbConfig)
	case "sqlite":
		return CreateSqlite(dbConfig)
	default:
		return fmt.Errorf("adapter no supported")
	}
}

func CreateMyslq(dbConfig *config.Database) (interface{}, error) {
	return sql.Open("mysql", config.Source)
}

func CreatePgsql(dbConfig *config.Database) (interface{}, error) {
	return sql.Open("postgres", config.Source)
}

func CreateSqlite(dbConfig *config.Database) (interface{}, error) {
	return sql.Open("sqlite3", config.Source)
}
*/