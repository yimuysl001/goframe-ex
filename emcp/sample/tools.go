package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/mark3labs/mcp-go/mcp"
	"goframe-ex/edb"
	"strings"
)

type LocalTableFiled struct{}

func (d *LocalTableFiled) New() mcp.Tool {
	return mcp.NewTool("getTableFiled",
		mcp.WithDescription("获取本地数据库相应字段"),
		mcp.WithArray("tables", mcp.Required(), mcp.Description("数据库名称")),
	)
}
func (d *LocalTableFiled) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	tables := gconv.Strings(args["tables"])
	if len(tables) == 0 {
		g.Log().Error(ctx, "入参获取失败")
		return nil, errors.New("tables 参数是必需的")
	}
	var tableMap = make(map[string][]*gdb.TableField)
	g.Log().Info(ctx, "===========LocalTableFiled============")
	for _, table := range tables {
		cchema := ""
		tableName := strings.Clone(strings.ToLower(table))
		if strings.Contains(tableName, "..") {
			splitN := strings.SplitN(tableName, "..", 2)
			cchema = splitN[0]
			tableName = splitN[1]
		}
		if strings.Contains(tableName, ".dbo.") {
			splitN := strings.SplitN(tableName, ".dbo.", 2)
			cchema = splitN[0]
			tableName = splitN[1]
		}

		fields, err := edb.DB("BS").GDB().TableFields(ctx, tableName, cchema)
		if err != nil {
			g.Log().Error(ctx, err)
			return nil, err
		}
		var ts = make([]*gdb.TableField, len(fields))
		var index = 0
		for _, field := range fields {
			ts[index] = field
			index++
		}
		tableMap[table] = ts

	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(gjson.New(tableMap).MustToJsonIndentString()),
		},
	}, nil

}

type ExplainTableName struct{}

func (e *ExplainTableName) New() mcp.Tool {
	return mcp.NewTool("explainTableName",
		mcp.WithDescription(`**🚀 根据配置将表结构说明导入到数据库**
重要：ExecutionPlan结构体格式要求：
[
    {
        "tableName": "表名称(string)",
        "fields": [
            {
                "fieldName": "字段名(string)",
                "fieldDesc": "字段描述(string)",
                "fieldType": "字段类型:string/int/bool/time.Time等(string)",
                "fieldRemark": "字段备注"
            }
        ]
    }
]
注意：
1. tableName如果相同，必须合并相应的数据
2. fieldName必须数据必须为英文字母或者英文字母加数字的组合，不能为中文
`),
		mcp.WithArray("executionPlan", mcp.Required(), mcp.Description("数据结构说明")),
	)
}

type ExecutionPlan struct {
	TableName string
	Fields    []*TableField
}

type TableField struct {
	TableName   string
	FieldName   string
	FieldDesc   string
	FieldType   string
	FieldRemark string
}

func (e *ExplainTableName) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	g.Log().Info(ctx, "===========ExplainTableName============")
	var index = 0
	args := request.GetArguments()
	executionPlans := gconv.Interfaces(args["executionPlan"])
	var tbf = make([]ExecutionPlan, 0)

	err := gconv.Scan(executionPlans, &tbf)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, err
	}
	if len(tbf) == 0 {
		g.Log().Error(ctx, "未生成数据")
		return nil, errors.New("未生成数据")
	}
	var model = edb.DB().GDB().Schema("YXJCPTJKK").Model("TB_TABLEDESC").Safe(true)
	for _, executionPlan := range tbf {
		if len(executionPlan.Fields) == 0 {
			continue
		}
		index = index + len(executionPlan.Fields)
		for _, field := range executionPlan.Fields {
			field.TableName = executionPlan.TableName
		}
		_, err = model.Where("TableName", executionPlan.TableName).Delete()
		if err != nil {
			g.Log().Error(ctx, err)
			return nil, err
		}
		_, err = model.Batch(20).Insert(executionPlan.Fields)
		if err != nil {
			g.Log().Error(ctx, err)
			return nil, err
		}

	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(fmt.Sprintf("生成完成，共生成%d条数据", index)),
		},
	}, nil
}
