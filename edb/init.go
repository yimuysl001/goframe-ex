package edb

import (
	_ "github.com/gogf/gf/contrib/drivers/clickhouse/v2" //加载数据数据库驱动
	_ "github.com/gogf/gf/contrib/drivers/dm/v2"         //加载数据数据库驱动
	_ "github.com/gogf/gf/contrib/drivers/mssql/v2"      //加载数据数据库驱动
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"      //加载数据数据库驱动
	_ "github.com/gogf/gf/contrib/drivers/oracle/v2"     //加载数据数据库驱动
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"      //加载数据数据库驱动
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"     //加载数据数据库驱动
	"goframe-ex/egoja/gojaapi"
)

func init() {

	gojaapi.RegisterCommonParameter("db", DB)
}
