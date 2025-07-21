package edb

import (
	"context"
	"strings"
	"time"
)

// 获取数据库所有库

const (
	databaseSqlMssql    = "select name AS database from sysdatabases"
	databaseSqlMysql    = "SELECT SCHEMA_NAME AS `database` FROM INFORMATION_SCHEMA.SCHEMATA"
	databaseSqlOracle   = `SELECT username AS "database" FROM all_users`
	databaseSqlPgsql    = `SELECT SCHEMA_NAME as database FROM pg_namespace `
	database_sql_sqlite = ``
)

const (
	nowSqlMssql    = "SELECT GETDATE() as cnow"
	nowSqlMysql    = "SELECT NOW() as `cnow`"
	nowSqlOracle   = `SELECT SYSTIMESTAMP as "cnow" FROM DUAL`
	nowSqlPgsql    = `SELECT NOW() as cnow`
	now_sql_sqlite = `SELECT DATETIME('now') as cnow`
)

func (d *DataBase) GetSchemas(ctx context.Context) ([]string, error) {
	var sqls = ""
	switch strings.ToLower(d.gdb.GetConfig().Type) {
	case "mysql":
		sqls = databaseSqlMysql
	case "oracle":
		sqls = databaseSqlOracle
	case "postgres", "pgsql":
		sqls = databaseSqlPgsql
	case "sqlserver", "mssql":
		sqls = databaseSqlMssql

	}
	if sqls == "" {
		return []string{d.gdb.GetSchema()}, nil
	}

	array, err := d.gdb.GetArray(ctx, sqls)
	if err != nil {
		return []string{d.gdb.GetSchema()}, nil
	}
	schemas := make([]string, len(array))
	for i, schema := range array {
		schemas[i] = schema.String()
	}

	return schemas, nil

}

// GetSchemas 获取当前连接所有库
func GetSchemas(ctx context.Context, name string) ([]string, error) {
	return DB(name).GetSchemas(ctx)
	//var sqls = ""
	//switch strings.ToLower(db.GetConfig().Type) {
	//case "mysql":
	//	sqls = databaseSqlMysql
	//case "oracle":
	//	sqls = databaseSqlOracle
	//case "postgres", "pgsql":
	//	sqls = databaseSqlPgsql
	//case "sqlserver", "mssql":
	//	sqls = databaseSqlMssql
	//
	//}
	//
	//if sqls == "" {
	//	return []string{db.GetSchema()}, nil
	//}
	//
	//array, err := db.GetArray(ctx, sqls)
	//if err != nil {
	//	return []string{db.GetSchema()}, nil
	//}
	//schemas := make([]string, len(array))
	//for i, schema := range array {
	//	schemas[i] = schema.String()
	//}
	//
	//return schemas, nil
}

func (d *DataBase) GetNow(ctx context.Context) (time.Time, error) {
	var sqls = ""
	switch strings.ToLower(d.gdb.GetConfig().Type) {
	case "mysql":
		sqls = nowSqlMysql
	case "oracle":
		sqls = nowSqlOracle
	case "postgres", "pgsql":
		sqls = nowSqlPgsql
	case "sqlserver", "mssql":
		sqls = nowSqlMssql
	default:
		return time.Now(), nil

	}
	value, err := d.gdb.GetValue(ctx, sqls)
	if err != nil || value == nil {
		return time.Now(), nil
	}
	return value.Time(), nil
}

func GetNow(ctx context.Context, name string) (time.Time, error) {
	return DB(name).GetNow(ctx)
	//db := DB(name).GDB()
	//
	//var sqls = ""
	//switch strings.ToLower(db.GetConfig().Type) {
	//case "mysql":
	//	sqls = nowSqlMysql
	//case "oracle":
	//	sqls = nowSqlOracle
	//case "postgres", "pgsql":
	//	sqls = nowSqlPgsql
	//case "sqlserver", "mssql":
	//	sqls = nowSqlMssql
	//default:
	//	return time.Now(), nil
	//
	//}
	//
	//value, err := db.GetValue(ctx, sqls)
	//if err != nil || value == nil {
	//	return time.Now(), nil
	//}
	//
	//return value.Time(), nil
}
