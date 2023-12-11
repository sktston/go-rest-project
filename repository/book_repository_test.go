package repository

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go-rest-project/config"
	"go-rest-project/model/entity"
	"go-rest-project/test"
	"os"
	"testing"
)

// Constants

func getTestBookA() entity.Book {
	return entity.Book{
		Title:     "TestTitleA",
		Author:    "TestAuthorA",
		Publisher: "TestPublisherA",
	}
}

func getTestBookB() entity.Book {
	return entity.Book{
		Title:     "TestTitleB",
		Author:    "TestAuthorB",
		Publisher: "TestPublisherB",
	}
}

// Tests

func TestCreateBook(t *testing.T) {
	// prepare
	test.InitTestDB(t)
	defer test.FreeTestDB(t)

	// test
	testBookA := getTestBookA()
	assert.NoError(t, CreateBook(&testBookA))
}

func TestGetBookList(t *testing.T) {
	// prepare
	test.InitTestDB(t)
	defer test.FreeTestDB(t)

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
	test.InitTestDB(t)
	defer test.FreeTestDB(t)

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
	test.InitTestDB(t)
	defer test.FreeTestDB(t)

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
	test.InitTestDB(t)
	defer test.FreeTestDB(t)

	testBookA := getTestBookA()
	assert.NoError(t, CreateBook(&testBookA))

	// test
	var book entity.Book
	assert.NoError(t, DeleteBook(&book, 1))
	assert.Error(t, GetBookByID(&book, 1))
}

// Helpers

// TestMain main function with postgres database
func TestMain(m *testing.M) {
	viper.Set("log.level", "TEST")
	config.SetLogLevel()

	// create postgres docker container
	pool, resource, err := test.CreatePostgres()
	if err != nil {
		fmt.Printf("Could not create postgres: %s", err)
		os.Exit(-1)
	}

	//Run tests
	code := m.Run()

	// remove postgres docker container
	if err := test.RemovePostgres(pool, resource); err != nil {
		fmt.Printf("Could not remove postgres: %s", err)
		os.Exit(-1)
	}

	os.Exit(code)
}
