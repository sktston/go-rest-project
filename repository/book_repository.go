package repository

import (
	"github.com/sktston/go-rest-project/db"
	"github.com/sktston/go-rest-project/model/entity"
)

//CreateBook ... Insert New data
func CreateBook(book *entity.Book) (err error) {
	if err = db.GetDB().Create(book).Error; err != nil {
		return err
	}
	return nil
}

//GetBookList Fetch all book data
func GetBookList(book *[]entity.Book) (err error) {
	if err = db.GetDB().Find(book).Error; err != nil {
		return err
	}
	return nil
}

//GetBookByID ... Fetch only one book by Id
func GetBookByID(book *entity.Book, id int) (err error) {
	if err = db.GetDB().First(book, id).Error; err != nil {
		return err
	}
	return nil
}

//UpdateBook ... Update book
func UpdateBook(book *entity.Book) (err error) {
	db.GetDB().Save(book)
	return nil
}

//DeleteBook ... Delete book
func DeleteBook(book *entity.Book, id int) (err error) {
	db.GetDB().Delete(book, id)
	return nil
}