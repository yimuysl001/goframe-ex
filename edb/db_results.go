package edb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/ghash"
	"github.com/gogf/gf/v2/util/gconv"
	"goframe-ex/egoja/eparser"
)

func (d *DataBase) Exec(ctx context.Context, sql string, params map[string]any) (result sql.Result, err error) {

	parseSql, parameters, err := eparser.ParseSql(sql, params)
	if err != nil {
		return nil, err
	}
	return d.gdb.Exec(ctx, parseSql, parameters...)

}

func (d *DataBase) SelectOne(ctx context.Context, sql string, params map[string]any) (result gdb.Record, err error) {
	var ckey = d.cacheKey
	if d.cachetime > 0 {
		if ckey == "" || ckey == fmt.Sprintf(`SelectCache:%s:`, d.name) {
			ckey = fmt.Sprintf(`SelectCache:%s:%d`, d.name, ghash.BKDR64([]byte(sql+", @PARAMS:"+gconv.String(params))))
		}
		get, err := d.gdb.GetCache().Get(ctx, ctx)

		if err == nil && get != nil {
			err = get.Scan(&result)
			return result, err
		}
	}
	parseSql, parameters, err := eparser.ParseSql(sql, params)
	if err != nil {
		return nil, err
	}

	exec, err := d.gdb.GetOne(ctx, parseSql, parameters...)
	if err != nil {
		return nil, err
	}
	if d.cachetime > 0 && (len(exec) > 0 || d.nullCache) {
		_ = d.gdb.GetCache().Set(ctx, ckey, exec, d.cachetime)
	}
	return exec, err
}

func (d *DataBase) Select(ctx context.Context, sql string, params map[string]any) (result gdb.Result, err error) {
	var ckey = d.cacheKey
	if d.cachetime > 0 {
		if ckey == "" || ckey == fmt.Sprintf(`SelectCache:%s:`, d.name) {
			ckey = fmt.Sprintf(`SelectCache:%s:%d`, d.name, ghash.BKDR64([]byte(sql+", @PARAMS:"+gconv.String(params))))
		}
		get, err := d.gdb.GetCache().Get(ctx, ctx)

		if err == nil && get != nil {
			err = get.Scan(&result)
			return result, err
		}
	}
	parseSql, parameters, err := eparser.ParseSql(sql, params)
	if err != nil {
		return nil, err
	}

	exec, err := d.gdb.Query(ctx, parseSql, parameters...)
	if err != nil {
		return nil, err
	}
	if d.cachetime > 0 && (len(exec) > 0 || d.nullCache) {
		_ = d.gdb.GetCache().Set(ctx, ckey, exec, d.cachetime)
	}
	return exec, err
}

func (d *DataBase) SelectList(ctx context.Context, sql string, params map[string]any) (result []gdb.Result, err error) {
	var ckey = d.cacheKey
	if d.cachetime > 0 {
		if ckey == "" || ckey == fmt.Sprintf(`SelectCache:%s:`, d.name) {
			ckey = fmt.Sprintf(`SelectCache:%s:%d`, d.name, ghash.BKDR64([]byte(sql+", @PARAMS:"+gconv.String(params))))
		}
		get, err := d.gdb.GetCache().Get(ctx, ctx)

		if err == nil && get != nil {
			err = get.Scan(&result)
			return result, err
		}
	}

	parseSql, parameters, err := eparser.ParseSql(sql, params)
	if err != nil {
		return nil, err
	}

	exec, err := d.GetResults(ctx, true, "", parseSql, parameters...)
	if err != nil {
		return nil, err
	}
	if d.cachetime > 0 && (len(exec) > 0 || d.nullCache) {
		_ = d.gdb.GetCache().Set(ctx, ckey, exec, d.cachetime)
	}
	return exec, err
}

func (d *DataBase) GetResults(ctx context.Context, master bool, schema string, sqlstr string, args ...interface{}) ([]gdb.Result, error) {

	link, err := d.gdb.GetCore().GetLink(ctx, master, schema)
	if err != nil {
		return nil, err
	}
	queryContext, err := link.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer queryContext.Close()
	var results = make([]gdb.Result, 0)

	if dt, err := getResult(queryContext); err != nil {
		return nil, err
	} else {
		// 将当前数据集对象加入到数据集列表对象
		results = append(results, dt)
	}
	for queryContext.NextResultSet() {
		if dt, err := getResult(queryContext); err != nil {
			return nil, err
		} else {
			//fmt.Println("处理下一个数据集")
			// 将当前数据集对象加入到数据集列表对象
			results = append(results, dt)
			// 更新数据集列表包含的数据集个数
		}
	}
	return results, nil
}

func getResult(rows *sql.Rows) (gdb.Result, error) {
	var result = make(gdb.Result, 0)
	if !rows.Next() {
		return result, nil
	}
	// Column names and types.
	columns, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	//columnTypes := make([]string, len(columns))
	columnNames := make([]string, len(columns))
	for k, v := range columns {
		//columnTypes[k] = v.DatabaseTypeName()
		columnNames[k] = v.Name()
	}
	var (
		values   = make([]interface{}, len(columnNames))
		scanArgs = make([]interface{}, len(values))
	)

	for i := range values {
		scanArgs[i] = &values[i]
	}
	for {
		if err := rows.Scan(scanArgs...); err != nil {
			return result, err
		}
		record := gdb.Record{}
		for i, value := range values {
			record[columnNames[i]] = gvar.New(value)
		}
		result = append(result, record)
		if !rows.Next() {
			break
		}
	}
	return result, nil
}
