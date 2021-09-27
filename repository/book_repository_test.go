package repository

import (
	"github.com/rs/zerolog/log"
	"github.com/sktston/go-rest-project/model/entity"
	"github.com/sktston/go-rest-project/test"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// Constants

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

// Tests

func TestCreateBook(t *testing.T) {
	// prepare
	testDB := test.InitTestDB(t)
	defer test.FreeTestDB(t, testDB)

	// test
	testBookA := getTestBookA()
	assert.NoError(t, CreateBook(&testBookA))
}

func TestGetBookList(t *testing.T) {
	// prepare
	testDB := test.InitTestDB(t)
	defer test.FreeTestDB(t, testDB)

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
	testDB := test.InitTestDB(t)
	defer test.FreeTestDB(t, testDB)

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
	testDB := test.InitTestDB(t)
	defer test.FreeTestDB(t, testDB)

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
	testDB := test.InitTestDB(t)
	defer test.FreeTestDB(t, testDB)

	testBookA := getTestBookA()
	assert.NoError(t, CreateBook(&testBookA))

	// test
	var book entity.Book
	assert.NoError(t, DeleteBook(&book, 1))
	assert.Error(t, GetBookByID(&book, 1))
}

// Helpers

func TestMain(m *testing.M) {
	// create postgres docker
	pool, resource, err := test.CreatePostgres()
	if err != nil {
		log.Fatal().Msgf("Could not create postgres: %s", err)
	}

	//Run tests
	code := m.Run()

	// delete postgres docker
	if err := test.DeletePostgres(pool, resource); err != nil {
		log.Fatal().Msgf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}