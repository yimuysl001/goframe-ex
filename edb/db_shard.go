package edb

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
)

type ShardingTableType int
type ShardingSchemaType int

const (
	DTABLE  = ShardingTableType(iota) // 单表  tableName
	MTABLE                            // 月表 tableName01
	YTABLE                            // 年表 tableName2024
	YMTABLE                           // 年月表 tableName202401
	FTABLE                            // 分区表 tableName_01
	DFTABLE                           //  自定义分区表 tableNameAAA
)
const (
	DSCHEMA  = ShardingSchemaType(iota) // 单库  SCHEMA
	MSCHEMA                             // 月库  SCHEMA01
	YSCHEMA                             // 年表 SCHEMA2024
	YMSCHEMA                            // 年月表 SCHEMA202401
	FSCHEMA                             // 分区表 SCHEMA_01
	DFSCHEMA                            // 自定义库 SCHEMAAA
)

type ShardingRule struct {
	TableType  ShardingTableType
	SCHEMAType ShardingSchemaType
	Key        any
}

func NewShardingConfig(tableType ShardingTableType, shardingSchemaType ShardingSchemaType, key any) ShardingRule {
	return ShardingRule{tableType, shardingSchemaType, key}
}

func (s *ShardingRule) ShardingConfig(PrefixTable, PrefixTableSchemas string, enable bool) gdb.ShardingConfig {
	return gdb.ShardingConfig{
		Table: gdb.ShardingTableConfig{
			Enable: enable,
			Prefix: PrefixTable,
			Rule:   s,
		},
		Schema: gdb.ShardingSchemaConfig{
			Enable: enable,
			Prefix: PrefixTableSchemas,
			Rule:   s,
		},
	}
}

func (s *ShardingRule) TableName(ctx context.Context, config gdb.ShardingTableConfig, value any) (string, error) {
	var datas = s.Key

	switch s.TableType {
	case DTABLE:
		return config.Prefix, nil
	case MTABLE:
		switch v := datas.(type) {
		case string:
			return config.Prefix + v, nil
		default:
			return config.Prefix + gconv.Time(datas, "2006-01-02 15:04:05").Format("02"), nil
		}
	case YTABLE:
		switch v := datas.(type) {
		case string:
			return config.Prefix + v, nil
		default:
			return config.Prefix + gconv.Time(datas, "2006-01-02 15:04:05").Format("2006"), nil
		}
	case YMTABLE:
		switch v := datas.(type) {
		case string:
			return config.Prefix + v, nil
		default:
			return config.Prefix + gconv.Time(datas, "2006-01-02 15:04:05").Format("200601"), nil
		}
	case FTABLE:
		var data = strings.TrimSpace(gconv.String(datas))
		if len(data) == 1 {
			return config.Prefix + "0" + data, nil
		} else {
			return config.Prefix + data[len(data)-2:], nil
		}

	case DFTABLE:
		return config.Prefix + gconv.String(datas), nil
	}
	return config.Prefix, nil
}

func (s *ShardingRule) SchemaName(ctx context.Context, config gdb.ShardingSchemaConfig, value any) (string, error) {
	var datas = s.Key

	switch s.SCHEMAType {
	case DSCHEMA:
		return config.Prefix, nil
	case MSCHEMA:
		switch v := datas.(type) {
		case string:
			return config.Prefix + v, nil
		default:
			return config.Prefix + gconv.Time(datas, "2006-01-02 15:04:05").Format("02"), nil
		}
	case YSCHEMA:
		switch v := datas.(type) {
		case string:
			return config.Prefix + v, nil
		default:
			return config.Prefix + gconv.Time(datas, "2006-01-02 15:04:05").Format("2006"), nil
		}
	case YMSCHEMA:
		switch v := datas.(type) {
		case string:
			return config.Prefix + v, nil
		default:
			return config.Prefix + gconv.Time(datas, "2006-01-02 15:04:05").Format("200601"), nil
		}
	case FSCHEMA:
		var data = strings.TrimSpace(gconv.String(datas))
		if len(data) == 1 {
			return config.Prefix + "0" + data, nil
		} else {
			return config.Prefix + data[len(data)-2:], nil
		}

	case DFSCHEMA:
		return config.Prefix + gconv.String(datas), nil
	}
	return config.Prefix, nil

}
