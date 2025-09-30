package main

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/xuri/excelize/v2"
	"goframe-ex/edb"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type SQLServerXLSXImporter struct {
	ctx context.Context
	db  gdb.DB
}

func NewSQLServerXLSXImporter(name string) (*SQLServerXLSXImporter, error) {

	return &SQLServerXLSXImporter{ctx: gctx.New(), db: edb.DB(name).GDB()}, nil
}

// 推断 SQL Server 数据类型
func inferSQLServerColumnType(values []string) string {
	if len(values) == 0 {
		return "NVARCHAR(500)"
	}

	isInt := true
	isFloat := true
	isBit := true
	isDate := true
	isDateTime := true

	for _, value := range values {
		if value == "" {
			continue
		}

		value = strings.TrimSpace(value)

		// 检查整数
		if _, err := strconv.Atoi(value); err != nil {
			isInt = false
		}

		// 检查浮点数
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			isFloat = false
		}

		// 检查布尔值
		if value != "0" && value != "1" &&
			strings.ToLower(value) != "true" &&
			strings.ToLower(value) != "false" &&
			strings.ToLower(value) != "yes" &&
			strings.ToLower(value) != "no" {
			isBit = false
		}

		// 检查日期
		if isDate {
			if _, err := time.Parse("2006-01-02", value); err != nil {
				if _, err := time.Parse("2006/01/02", value); err != nil {
					isDate = false
				}
			}
		}

		// 检查日期时间
		if isDateTime {
			if _, err := time.Parse("2006-01-02 15:04:05", value); err != nil {
				if _, err := time.Parse("2006/01/02 15:04:05", value); err != nil {
					if _, err := time.Parse(time.RFC3339, value); err != nil {
						isDateTime = false
					}
				}
			}
		}
	}

	// 数据类型优先级
	if isBit && len(values) > 0 {
		// 只有全部是布尔值时才使用 BIT
		allBoolean := true
		for _, value := range values {
			lowerVal := strings.ToLower(value)
			if value != "0" && value != "1" &&
				lowerVal != "true" && lowerVal != "false" &&
				lowerVal != "yes" && lowerVal != "no" {
				allBoolean = false
				break
			}
		}
		if allBoolean {
			return "BIT"
		}
	}

	if isInt {
		return "NVARCHAR(255)"
	} else if isFloat {
		return "DECIMAL(18,2)"
	} else if isDateTime {
		return "DATETIME"
	} else if isDate {
		return "DATE"
	} else {
		// 根据最大长度调整 NVARCHAR 大小
		maxLen := 0
		for _, value := range values {
			if len(value) > maxLen {
				maxLen = len(value)
			}
		}

		if maxLen == 0 {
			return "NVARCHAR(500)"
		} else if maxLen <= 2000 {
			return "NVARCHAR(2000)"
		} else if maxLen <= 4000 {
			return "NVARCHAR(4000)"
		} else {
			return "NVARCHAR(MAX)"
		}
	}
}

// 生成安全的 SQL Server 列名
func safeSQLServerColumnName(name string) string {
	// 移除特殊字符，只保留字母、数字和下划线
	reg := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	safe := reg.ReplaceAllString(name, "")

	if safe == "" {
		return "column_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	// 确保不以数字开头
	if matched, _ := regexp.MatchString(`^[0-9]`, safe); matched {
		safe = "col_" + safe
	}

	// SQL Server 关键字检查
	sqlServerKeywords := map[string]bool{
		"key": true, "table": true, "select": true, "insert": true,
		"update": true, "delete": true, "from": true, "where": true,
		"group": true, "order": true, "by": true, "join": true,
		"left": true, "right": true, "inner": true, "outer": true,
	}

	if sqlServerKeywords[strings.ToLower(safe)] {
		safe = "col_" + safe
	}

	return safe
}

// 检查表是否存在
func (si *SQLServerXLSXImporter) tableExists(tableName string) (bool, error) {

	query := `
        SELECT CASE WHEN EXISTS (
            SELECT 1 FROM INFORMATION_SCHEMA.TABLES 
            WHERE TABLE_NAME = ?
        ) THEN 1 ELSE 0 END as a
    `

	value, err2 := si.db.GetValue(si.ctx, query, tableName)
	if err2 != nil {
		return false, err2
	}

	return value.Bool(), err2
}

// 创建表（SQL Server 语法）
func (si *SQLServerXLSXImporter) createTable(tableName string, headers []string, sampleData [][]string) error {
	if len(sampleData) == 0 {
		return fmt.Errorf("no sample data provided")
	}

	// 检查表是否存在，如果存在则删除
	exists, err := si.tableExists(tableName)
	if err != nil {
		return err
	}

	if exists {
		log.Printf("Table %s already exists, dropping...", tableName)
		dropSQL := fmt.Sprintf("DROP TABLE [%s]", tableName)
		_, err = si.db.Exec(si.ctx, dropSQL)
		if err != nil {
			return err
		}
	}

	var columns []string
	for i, header := range headers {
		if i >= len(sampleData[0]) {
			continue
		}

		// 收集该列的所有样本值
		var columnValues []string
		for _, row := range sampleData {
			if i < len(row) {
				columnValues = append(columnValues, row[i])
			}
		}

		colType := inferSQLServerColumnType(columnValues)
		colName := safeSQLServerColumnName(header)
		columns = append(columns, fmt.Sprintf("[%s] %s NULL", colName, colType))
	}

	createSQL := fmt.Sprintf(`
        CREATE TABLE [%s] (
            %s,
            [created_at] DATETIME DEFAULT GETDATE()
        )`, tableName, strings.Join(columns, ",\n"))

	_, err = si.db.Exec(si.ctx, createSQL)
	return err
}

// 批量插入数据到 SQL Server
func (si *SQLServerXLSXImporter) batchInsertData(tableName string, headers []string, data [][]string, batchSize int) error {
	// 准备列名
	safeHeaders := make([]string, len(headers))
	for i, header := range headers {
		safeHeaders[i] = safeSQLServerColumnName(header)
	}

	// 构建基础插入语句
	columnsStr := strings.Join(safeHeaders, ",")
	placeholders := "?" + strings.Repeat(",?", len(safeHeaders)-1)

	baseSQL := fmt.Sprintf("INSERT INTO [%s] (%s) VALUES (%s)",
		tableName, columnsStr, placeholders)

	totalRows := len(data)
	//var valueArgs []interface{}

	// 开始事务
	tx, err := si.db.Begin(si.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback() // 确保在出错时回滚

	//stmt, err := tx.Prepare(baseSQL)
	//if err != nil {
	//	return err
	//}
	//defer stmt.Close()

	for i, row := range data {
		// 确保行数据长度与表头一致
		paddedRow := make([]string, len(headers))
		copy(paddedRow, row)

		// 清空当前批次的参数
		currentArgs := make([]interface{}, len(headers))
		for j, value := range paddedRow {
			// 处理空值
			if value == "" {
				currentArgs[j] = nil
			} else {
				currentArgs[j] = value
			}
		}

		// 执行插入
		_, err := tx.Exec(baseSQL, currentArgs...)
		if err != nil {
			log.Printf("Error inserting row %d: %v", i+1, err)
			return err
		}

		// 每 batchSize 行提交一次事务并开始新事务
		if (i+1)%batchSize == 0 {
			err = tx.Commit()
			if err != nil {
				return err
			}

			log.Printf("Committed batch: %d/%d rows", i+1, totalRows)

			// 开始新事务
			tx, err = si.db.Begin(si.ctx)
			if err != nil {
				return err
			}

			//stmt, err = tx.Prepare(baseSQL)
			//if err != nil {
			//	return err
			//}
		}

		// 更新进度
		if (i+1)%1000 == 0 || i == totalRows-1 {
			log.Printf("Processed %d/%d rows (%.2f%%)",
				i+1, totalRows, float64(i+1)/float64(totalRows)*100)
		}
	}

	// 提交最后一批数据
	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Printf("Successfully imported %d rows into table [%s]", totalRows, tableName)
	return nil
}

// 主导入函数
func (si *SQLServerXLSXImporter) ImportXLSX(filePath, tableName string, batchSize int) error {
	// 打开 Excel 文件
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// 获取第一个工作表
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return err
	}

	if len(rows) < 2 {
		return fmt.Errorf("insufficient data in Excel file")
	}

	// 提取表头
	headers := rows[0]
	log.Printf("Found %d columns: %v", len(headers), headers)

	// 提取样本数据用于推断类型（前100行）
	sampleSize := min(100, len(rows)-1)
	sampleData := make([][]string, sampleSize)
	for i := 0; i < sampleSize; i++ {
		sampleData[i] = rows[i+1]
	}

	log.Printf("Using %d sample rows for type inference", len(sampleData))

	// 创建表
	if err := si.createTable(tableName, headers, sampleData); err != nil {
		return err
	}

	log.Printf("Table [%s] created successfully", tableName)

	// 批量插入数据
	return si.batchInsertData(tableName, headers, rows[1:], batchSize)
}

// 流式导入大文件（内存优化）
func (si *SQLServerXLSXImporter) StreamImportXLSX(filePath, tableName string, batchSize int) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)

	// 流式读取行
	rows, err := f.Rows(sheetName)
	if err != nil {
		return err
	}

	var headers []string
	var sampleData [][]string
	var allData [][]string

	rowIndex := 0
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return err
		}

		if rowIndex == 0 {
			headers = row
			log.Printf("Found %d columns: %v", len(headers), headers)
		} else {
			if rowIndex <= 100 {
				sampleData = append(sampleData, row)
			}
			allData = append(allData, row)
		}
		rowIndex++

		// 进度显示
		if rowIndex%1000 == 0 {
			log.Printf("Read %d rows from Excel", rowIndex)
		}
	}

	if err := si.createTable(tableName, headers, sampleData); err != nil {
		return err
	}

	log.Printf("Table [%s] created successfully, starting data import...", tableName)
	return si.batchInsertData(tableName, headers, allData, batchSize)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
