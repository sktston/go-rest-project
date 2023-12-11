package repository

import (
	"go-rest-project/database"
	"go-rest-project/model/entity"
)

// CreateBook ... Insert New data
func CreateBook(book *entity.Book) (err error) {
	if err = database.GetDB().Create(book).Error; err != nil {
		return err
	}
	return nil
}

// GetBookList Fetch all book data
func GetBookList(book *[]entity.Book) (err error) {
	if err = database.GetDB().Find(book).Error; err != nil {
		return err
	}
	return nil
}

// GetBookByID ... Fetch only one book by Id
func GetBookByID(book *entity.Book, id int) (err error) {
	if err = database.GetDB().First(book, id).Error; err != nil {
		return err
	}
	return nil
}

// UpdateBook ... Update book
func UpdateBook(book *entity.Book) (err error) {
	database.GetDB().Save(book)
	return nil
}

// DeleteBook ... Delete book
func DeleteBook(book *entity.Book, id int) (err error) {
	database.GetDB().Delete(book, id)
	return nil
}
