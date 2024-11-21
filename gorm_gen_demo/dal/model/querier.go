package model

import "gorm.io/gen"

// 通过添加注释生成自定义方法

type Querier interface {
	// SELECT * FROM @@table WHERE id = @id
	GetByID(id int) (gen.T, error)

	// GetByIDReturnMap 根据ID查询返回map
	//
	// SELECT * FROM @@table WHERE id=@id
	GetByIDReturnMap(id int) (gen.M, error) // 返回 map 和 error

	// SELECT * FROM @@table WHERE author=@author
	GetBooksByAuthor(author string) ([]*gen.T, error) // 返回数据切片和 error
}

type Searcher interface {
	// Search 根据指定条件查询书籍
	//
	// SELECT * FROM book
	// WHERE publish_date is not null
	// {{if book != nil}}
	//   {{if book.ID > 0}}
	//     AND id = @book.ID
	//   {{else if book.Author != ""}}
	//     AND author=@book.Author
	//   {{end}}
	// {{end}}
	Search(book *gen.T) ([]*gen.T, error)
}
