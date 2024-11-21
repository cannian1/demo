package main

import (
	"context"
	"demo/gorm_gen_demo/dal/model"
	"demo/gorm_gen_demo/dal/query"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

const dsn = "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True"

func initDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(fmt.Errorf("connect db fail: %w", err))
	}
	return db
}

func main() {
	fmt.Println("gen demo start...")

	db := initDB(dsn)

	// 初始化
	query.SetDefault(db)

	// CRUD
	// 新增
	b1 := &model.Book{
		Title:       "Getting Rich",
		Author:      "WALLACE D.WATTLES",
		Price:       17,
		PublishDate: time.Date(1971, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	err := query.Book.WithContext(context.Background()).Create(b1)
	if err != nil {
		fmt.Println("创建书籍失败，err:", err)
		return
	}
	fmt.Println(b1)

	// 查询
	book := query.Book
	b, err := book.WithContext(context.Background()).First()
	// 也可以使用全局Q对象查询
	//book, err := query.Q.Book.WithContext(context.Background()).First()
	if err != nil {
		fmt.Printf("查询书籍失败, err:%v\n", err)
		return
	}
	fmt.Printf("book:%v\n", b)

	// 更新
	ret, err := book.
		WithContext(context.Background()).
		Where(book.ID.Eq(1)).
		Update(book.Price, 20)
	if err != nil {
		fmt.Println("更新书籍失败", err)
		return
	}
	fmt.Println(ret.RowsAffected)

	// 删除
	ret, err = query.Book.WithContext(context.Background()).Where(query.Book.ID.Eq(7)).Delete()
	if err != nil {
		fmt.Printf("删除书籍失败, err:%v\n", err)
		return
	}
	fmt.Printf("RowsAffected:%v\n", ret.RowsAffected)

	b2 := &model.Book{ID: 1}
	ret, err = book.WithContext(context.Background()).Delete(b2)
	if err != nil {
		fmt.Printf("删除书籍失败, err:%v\n", err)
		return
	}
	fmt.Printf("RowsAffected:%v\n", ret.RowsAffected)

	/*	// query.Q 用于嵌入结构体
		r := repo{db: *query.Q}*/

	returnMap, err := query.Book.WithContext(context.Background()).GetByIDReturnMap(2)
	if err != nil {
		return
	}
	fmt.Println(returnMap)

	//query.Book.WithContext(context.Background()).Search()
}

type repo struct {
	db query.Query
}
