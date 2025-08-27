package driver

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	//_ "github.com/marcboeker/go-duckdb/v2"
	_ "goframe-ex/edb/driver/duckdb"
	"log"
	"testing"
)

func TestDb(t *testing.T) {

	var ctx = gctx.New()
	fmt.Println(g.DB().Tables(ctx))

	//one, err := g.DB().GetOne(ctx, "select 1 as a ")
	//var model = g.DB().Model("products").Ctx(ctx).Safe(true)
	//exec, err := model.Insert(g.Map{
	//	"id":       2,
	//	"category": "Laptop",
	//	"name":     "Laptop",
	//	"price":    999.99,
	//})
	//
	//fmt.Println(exec, err)
	//fmt.Println(model.All())

	all, err := g.DB().GetAll(ctx, `SELECT
	*
FROM
	pg_class c
INNER JOIN pg_namespace n ON
	c.relnamespace = n.oid
WHERE

	 c.relkind IN ('r', 'p')`)

	fmt.Println(all)
	fmt.Println(err)

}

func TestDuck(t *testing.T) {
	// 连接数据库
	db, err := sql.Open("duckdb", "")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建表
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS products (
            id INTEGER PRIMARY KEY,
            name VARCHAR(100),
            price DECIMAL(10,2),
            category VARCHAR(50)
        )
    `)
	if err != nil {
		log.Fatal(err)
	}

	// 插入数据g
	_, err = db.Exec(`
        INSERT INTO products VALUES
        (1, 'Laptop', 999.99, 'Electronics'),
        (2, 'Book', 19.99, 'Education'),
        (3, 'Headphones', 49.99, 'Electronics')
    `)
	if err != nil {
		log.Fatal(err)
	}

	// 查询数据
	rows, err := db.Query(`
        SELECT category, COUNT(*), AVG(price) 
        FROM products 
        GROUP BY category
    `)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Category\tCount\tAvg Price")
	fmt.Println("-----------------------------------")
	for rows.Next() {
		var category string
		var count int
		var avgPrice float64

		err = rows.Scan(&category, &count, &avgPrice)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%-12s\t%d\t$%.2f\n", category, count, avgPrice)
	}
}
