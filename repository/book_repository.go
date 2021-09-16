package repository

import (
	"github.com/sktston/go-rest-project/config"
	"github.com/sktston/go-rest-project/models/entity"
)

//CreateBook ... Insert New data
func CreateBook(book *entity.Book) (err error) {
	if err = config.DB.Create(book).Error; err != nil {
		return err
	}
	return nil
}

//GetBookList Fetch all book data
func GetBookList(book *[]entity.Book) (err error) {
	if err = config.DB.Find(book).Error; err != nil {
		return err
	}
	return nil
}

//GetBookByID ... Fetch only one book by Id
func GetBookByID(book *entity.Book, id int) (err error) {
	if err = config.DB.First(book, id).Error; err != nil {
		return err
	}
	return nil
}

//UpdateBook ... Update book
func UpdateBook(book *entity.Book) (err error) {
	config.DB.Save(book)
	return nil
}

//DeleteBook ... Delete book
func DeleteBook(book *entity.Book, id int) (err error) {
	config.DB.Delete(book, id)
	return nil
}