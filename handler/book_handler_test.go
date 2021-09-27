package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sktston/go-rest-project/config"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
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
	testDB := config.InitTestDB(t)
	defer config.FreeTestDB(t, testDB)

	// test
	body, code := sendRequest(
		http.MethodPost,
		"/books",
		strings.NewReader(testBodyA),
		setupRouter(http.MethodPost, "/books", CreateBook),
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
	testDB := config.InitTestDB(t)
	defer config.FreeTestDB(t, testDB)

	assert.NoError(t, createBookA())
	assert.NoError(t, createBookB())

	// test
	body, code := sendRequest(
		http.MethodGet,
		"/books",
		nil,
		setupRouter(http.MethodGet, "/books", GetBookList),
	)
	assert.Equal(t, http.StatusOK, code)

	var res []map[string]interface{}
	assert.NoError(t, json.Unmarshal(body.Bytes(), &res))
	assert.Len(t, res, 2)
}

func TestGetBookByID(t *testing.T) {
	// prepare
	testDB := config.InitTestDB(t)
	defer config.FreeTestDB(t, testDB)

	assert.NoError(t, createBookA())

	// test
	body, code := sendRequest(
		http.MethodGet,
		"/books/1",
		nil,
		setupRouter(http.MethodGet, "/books/:id", GetBookByID),
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
	testDB := config.InitTestDB(t)
	defer config.FreeTestDB(t, testDB)

	assert.NoError(t, createBookA())

	// test
	updateBody := `{
		"title": "UpdatedTestTitleA",
		"author": "UpdatedTestAuthorA",
		"publisher": "UpdatedTestPublisherA"
	}`
	body, code := sendRequest(
		http.MethodPut,
		"/books/1",
		strings.NewReader(updateBody),
		setupRouter(http.MethodPut, "/books/:id", UpdateBook),
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
	testDB := config.InitTestDB(t)
	defer config.FreeTestDB(t, testDB)

	assert.NoError(t, createBookA())

	// test
	_, code := sendRequest(
		http.MethodDelete,
		"/books/1",
		nil,
		setupRouter(http.MethodDelete, "/books/:id", DeleteBook),
	)
	assert.Equal(t, http.StatusOK, code)

	_, code = sendRequest(
		http.MethodGet,
		"/books/1",
		nil,
		setupRouter(http.MethodGet, "/books/:id", GetBookByID),
	)
	assert.Equal(t, http.StatusNotFound, code)
}

// Helpers

// createBookA create book with testBodyA
func createBookA() error {
	_, code := sendRequest(
		http.MethodPost,
		"/books",
		strings.NewReader(testBodyA),
		setupRouter(http.MethodPost, "/books", CreateBook),
	)
	if code == http.StatusOK {
		return nil
	} else {
		return errors.New("createBookA failed")
	}
}

// createBookB create book with testBodyB
func createBookB() error {
	_, code :=  sendRequest(
		http.MethodPost,
		"/books",
		strings.NewReader(testBodyB),
		setupRouter(http.MethodPost, "/books", CreateBook),
	)
	if code == http.StatusOK {
		return nil
	} else {
		return errors.New("createBookB failed")
	}
}

// setupRouter get router on given handler
func setupRouter(httpMethod string, path string, handler gin.HandlerFunc) *gin.Engine {
	// prepare router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Handle(httpMethod, path, handler)
	return router
}

// sendRequest reads response from given http request.
func sendRequest(httpMethod string, target string, requestBody io.Reader, router *gin.Engine) (*bytes.Buffer, int) {
	// serve http on given response and request
	req := httptest.NewRequest(httpMethod, target, requestBody)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr.Body, rr.Code
}
