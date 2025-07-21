package edb

import (
	"context"
	"database/sql"
)

func (d *DataBase) GetSqlFields(ctx context.Context, master bool, schema string, sqlstr string, args ...interface{}) ([]*sql.ColumnType, error) {

	link, err := d.gdb.GetCore().GetLink(ctx, master, schema)
	if err != nil {
		return nil, err
	}
	queryContext, err := link.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer queryContext.Close()

	return queryContext.ColumnTypes()

}
