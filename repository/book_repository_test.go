package repository

import (
	"github.com/sktston/go-rest-project/config"
	"github.com/sktston/go-rest-project/model/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func enter(t *testing.T) *gorm.DB {
	assert.NoError(t, config.LoadConfig())
	testDB, err := config.InitTestDB()
	assert.NoError(t, err)
	return testDB
}

func leave(t *testing.T, testDB *gorm.DB) {
	assert.NoError(t, config.FreeTestDB(testDB))
}

func getTestBookList() []entity.Book {
	return []entity.Book{
		{Title: "TestTitleA", Author: "TestAuthorA", Publisher: "TestPublisherA" },
		{Title: "TestTitleB", Author: "TestAuthorB", Publisher: "TestPublisherB" },
	}
}

func TestCreateBook(t *testing.T) {
	testDB := enter(t)
	defer leave(t, testDB)

	for _, testBook := range getTestBookList() {
		err := CreateBook(&testBook)
		assert.NoError(t, err)
	}
}

func TestGetAllBooks(t *testing.T) {
	testDB := enter(t)
	defer leave(t, testDB)

	for _, testBook := range getTestBookList() {
		err := CreateBook(&testBook)
		assert.NoError(t, err)
	}

	var books []entity.Book
	err := GetBookList(&books)
	assert.NoError(t, err)

	assert.Len(t, books, len(getTestBookList()))

	for i, book := range books {
		assert.Equal(t, getTestBookList()[i].Title, book.Title)
		assert.Equal(t, getTestBookList()[i].Author, book.Author)
		assert.Equal(t, getTestBookList()[i].Publisher, book.Publisher)
	}
}

func TestGetBookByID(t *testing.T) {
	testDB := enter(t)
	defer leave(t, testDB)

	for _, testBook := range getTestBookList() {
		err := CreateBook(&testBook)
		assert.NoError(t, err)
	}

	var book entity.Book
	err := GetBookByID(&book, 1)
	assert.NoError(t, err)

	assert.Equal(t, getTestBookList()[0].Title, book.Title)
	assert.Equal(t, getTestBookList()[0].Author, book.Author)
	assert.Equal(t, getTestBookList()[0].Publisher, book.Publisher)
}

func TestUpdateBook(t *testing.T) {
	testDB := enter(t)
	defer leave(t, testDB)

	testBookC := entity.Book{Title: "TestTitleC", Author: "TestAuthorC", Publisher: "TestPublisherC" }

	for _, testBook := range getTestBookList() {
		err := CreateBook(&testBook)
		assert.NoError(t, err)
	}

	var book entity.Book
	err := GetBookByID(&book, 1)
	assert.NoError(t, err)

	book.Title = testBookC.Title
	book.Author = testBookC.Author
	book.Publisher = testBookC.Publisher

	err = UpdateBook(&book)
	assert.NoError(t, err)

	err = GetBookByID(&book, 1)
	assert.NoError(t, err)

	assert.Equal(t, testBookC.Title, book.Title)
	assert.Equal(t, testBookC.Author, book.Author)
	assert.Equal(t, testBookC.Publisher, book.Publisher)
}

func TestDeleteBook(t *testing.T) {
	testDB := enter(t)
	defer leave(t, testDB)

	for _, testBook := range getTestBookList() {
		err := CreateBook(&testBook)
		assert.NoError(t, err)
	}

	var book entity.Book
	err := DeleteBook(&book, 1)
	assert.NoError(t, err)

	err = GetBookByID(&book, 1)
	assert.Error(t, err)

	var books []entity.Book
	err = GetBookList(&books)
	assert.NoError(t, err)

	assert.Len(t, books, len(getTestBookList())-1)
}