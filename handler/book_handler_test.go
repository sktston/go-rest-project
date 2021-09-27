package handler

import (
	"encoding/json"
	"errors"
	"github.com/sktston/go-rest-project/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

// Constants

const (
	testBodyA = `{
		"title": "TestTitleA",
		"author": "TestAuthorA",
		"publisher": "TestPublisherA"
	}`
	testBodyB = `{
		"title": "TestTitleA",
		"author": "TestAuthorA",
		"publisher": "TestPublisherA"
	}`
)

// Tests

func TestCreateBook(t *testing.T) {
	// prepare
	testDB := test.InitTestDB(t)
	defer test.FreeTestDB(t, testDB)

	// test
	body, code := test.SendRequest(
		http.MethodPost,
		"/books",
		strings.NewReader(testBodyA),
		test.SetupRouter(http.MethodPost, "/books", CreateBook),
	)
	assert.Equal(t, http.StatusOK, code)

	res := make(map[string]interface{})
	assert.NoError(t, json.Unmarshal(body.Bytes(), &res))
	assert.Equal(t, 1, int(res["id"].(float64)))
	assert.Equal(t, "TestTitleA", res["title"].(string))
	assert.Equal(t, "TestAuthorA", res["author"].(string))
	assert.Equal(t, "TestPublisherA", res["publisher"].(string))
}

func TestGetBookList(t *testing.T) {
	// prepare
	testDB := test.InitTestDB(t)
	defer test.FreeTestDB(t, testDB)

	assert.NoError(t, createBookA())
	assert.NoError(t, createBookB())

	// test
	body, code := test.SendRequest(
		http.MethodGet,
		"/books",
		nil,
		test.SetupRouter(http.MethodGet, "/books", GetBookList),
	)
	assert.Equal(t, http.StatusOK, code)

	var res []map[string]interface{}
	assert.NoError(t, json.Unmarshal(body.Bytes(), &res))
	assert.Len(t, res, 2)
}

func TestGetBookByID(t *testing.T) {
	// prepare
	testDB := test.InitTestDB(t)
	defer test.FreeTestDB(t, testDB)

	assert.NoError(t, createBookA())

	// test
	body, code := test.SendRequest(
		http.MethodGet,
		"/books/1",
		nil,
		test.SetupRouter(http.MethodGet, "/books/:id", GetBookByID),
	)
	assert.Equal(t, http.StatusOK, code)

	res := make(map[string]interface{})
	assert.NoError(t, json.Unmarshal(body.Bytes(), &res))
	assert.Equal(t, 1, int(res["id"].(float64)))
	assert.Equal(t, "TestTitleA", res["title"].(string))
	assert.Equal(t, "TestAuthorA", res["author"].(string))
	assert.Equal(t, "TestPublisherA", res["publisher"].(string))
}

func TestUpdateBook(t *testing.T) {
	// prepare
	testDB := test.InitTestDB(t)
	defer test.FreeTestDB(t, testDB)

	assert.NoError(t, createBookA())

	// test
	updateBody := `{
		"title": "UpdatedTestTitleA",
		"author": "UpdatedTestAuthorA",
		"publisher": "UpdatedTestPublisherA"
	}`
	body, code := test.SendRequest(
		http.MethodPut,
		"/books/1",
		strings.NewReader(updateBody),
		test.SetupRouter(http.MethodPut, "/books/:id", UpdateBook),
	)
	assert.Equal(t, http.StatusOK, code)

	res := make(map[string]interface{})
	assert.NoError(t, json.Unmarshal(body.Bytes(), &res))
	assert.Equal(t, 1, int(res["id"].(float64)))
	assert.Equal(t, "UpdatedTestTitleA", res["title"].(string))
	assert.Equal(t, "UpdatedTestAuthorA", res["author"].(string))
	assert.Equal(t, "UpdatedTestPublisherA", res["publisher"].(string))
}

func TestDeleteBook(t *testing.T) {
	// prepare
	testDB := test.InitTestDB(t)
	defer test.FreeTestDB(t, testDB)

	assert.NoError(t, createBookA())

	// test
	_, code := test.SendRequest(
		http.MethodDelete,
		"/books/1",
		nil,
		test.SetupRouter(http.MethodDelete, "/books/:id", DeleteBook),
	)
	assert.Equal(t, http.StatusOK, code)

	_, code = test.SendRequest(
		http.MethodGet,
		"/books/1",
		nil,
		test.SetupRouter(http.MethodGet, "/books/:id", GetBookByID),
	)
	assert.Equal(t, http.StatusNotFound, code)
}

// Helpers

// createBookA create book with testBodyA
func createBookA() error {
	_, code := test.SendRequest(
		http.MethodPost,
		"/books",
		strings.NewReader(testBodyA),
		test.SetupRouter(http.MethodPost, "/books", CreateBook),
	)
	if code == http.StatusOK {
		return nil
	} else {
		return errors.New("createBookA failed")
	}
}

// createBookB create book with testBodyB
func createBookB() error {
	_, code :=  test.SendRequest(
		http.MethodPost,
		"/books",
		strings.NewReader(testBodyB),
		test.SetupRouter(http.MethodPost, "/books", CreateBook),
	)
	if code == http.StatusOK {
		return nil
	} else {
		return errors.New("createBookB failed")
	}
}