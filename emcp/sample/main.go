package main

import "log"

//func main() {
//	emcp.RegisterTool(&LocalTableFiled{})
//	emcp.RegisterTool(&ExplainTableName{})
//
//	s := g.Server()
//
//	emcp.RegisterConfig(s, emcp.DefaultConfig())
//	s.Run()
//
//}

func main() {
	// 初始化 SQL Server 导入器
	importer, err := NewSQLServerXLSXImporter("")
	if err != nil {
		log.Fatal("Failed to connect to SQL Server:", err)
	}

	log.Println("Connected to SQL Server successfully")

	// 方法1：普通导入（适合中小文件）
	err = importer.ImportXLSX("E:\\data\\xwechat_files\\wxid_rj8zox9d5ol222_5f9f\\msg\\attach\\fcbefc3b195771a2e08250647a0f87a2\\2025-09\\Rec\\acd627d08353ffc9\\F\\1\\2024国家药品信息标准库V11.0发布.xlsx", "TBZDYPXXV11", 1000)
	if err != nil {
		log.Fatal("Import failed:", err)
	}

	// 方法2：流式导入（适合大文件）
	// err = importer.StreamImportXLSX("large_data.xlsx", "large_imported_data", 500)
	// if err != nil {
	//     log.Fatal("Stream import failed:", err)
	// }

	log.Println("Import completed successfully!")
}
