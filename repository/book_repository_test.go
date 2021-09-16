package repository

import (
	"github.com/sktston/go-rest-project/config"
	"github.com/sktston/go-rest-project/model/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func getTestBookA() entity.Book {
	return entity.Book{
		Title: "TestTitleA",
		Author: "TestAuthorA",
		Publisher: "TestPublisherA",
	}
}

func getTestBookB() entity.Book {
	return entity.Book{
		Title: "TestTitleB",
		Author: "TestAuthorB",
		Publisher: "TestPublisherB",
	}
}

func TestCreateBook(t *testing.T) {
	// prepare
	testDB := startTestDB(t)
	defer finishTestDB(t, testDB)

	// test
	testBookA := getTestBookA()
	assert.NoError(t, CreateBook(&testBookA))
}

func TestGetBookList(t *testing.T) {
	// prepare
	testDB := startTestDB(t)
	defer finishTestDB(t, testDB)
	testBookA := getTestBookA()
	assert.NoError(t, CreateBook(&testBookA))
	testBookB := getTestBookB()
	assert.NoError(t, CreateBook(&testBookB))

	// test
	var books []entity.Book
	assert.NoError(t, GetBookList(&books))
	assert.Len(t, books, 2)
}

func TestGetBookByID(t *testing.T) {
	// prepare
	testDB := startTestDB(t)
	defer finishTestDB(t, testDB)

	testBookA := getTestBookA()
	assert.NoError(t, CreateBook(&testBookA))

	// test
	var book entity.Book
	assert.NoError(t, GetBookByID(&book, 1))

	assert.Equal(t, "TestTitleA", book.Title)
	assert.Equal(t, "TestAuthorA", book.Author)
	assert.Equal(t, "TestPublisherA", book.Publisher)
}

func TestUpdateBook(t *testing.T) {
	// prepare
	testDB := startTestDB(t)
	defer finishTestDB(t, testDB)

	testBookA := getTestBookA()
	assert.NoError(t, CreateBook(&testBookA))

	// test
	var book entity.Book
	assert.NoError(t, GetBookByID(&book, 1))

	book.Title = "UpdatedTestTitleA"
	book.Author = "UpdatedTestAuthorA"
	book.Publisher = "UpdatedTestPublisherA"
	assert.NoError(t, UpdateBook(&book))

	var updatedBook entity.Book
	assert.NoError(t, GetBookByID(&updatedBook, 1))

	assert.Equal(t, "UpdatedTestTitleA", updatedBook.Title)
	assert.Equal(t, "UpdatedTestAuthorA", updatedBook.Author)
	assert.Equal(t, "UpdatedTestPublisherA", updatedBook.Publisher)
}

func TestDeleteBook(t *testing.T) {
	// prepare
	testDB := startTestDB(t)
	defer finishTestDB(t, testDB)

	testBookA := getTestBookA()
	assert.NoError(t, CreateBook(&testBookA))

	// test
	var book entity.Book
	assert.NoError(t, DeleteBook(&book, 1))
	assert.Error(t, GetBookByID(&book, 1))
}

// startTestDB init test database
func startTestDB(t *testing.T) *gorm.DB {
	assert.NoError(t, config.LoadConfig())
	testDB, err := config.InitTestDB()
	assert.NoError(t, err)
	return testDB
}

// finishTestDB free test database
func finishTestDB(t *testing.T, testDB *gorm.DB) {
	assert.NoError(t, config.FreeTestDB(testDB))
}