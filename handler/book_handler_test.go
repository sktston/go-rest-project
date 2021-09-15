package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sktston/go-rest-project/config"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateBook(t *testing.T) {
	start(t)
	defer finish(t)

	// test
	testBody := `{
		"title": "TestTitleA",
		"author": "TestAuthorA",
		"publisher": "TestPublisherA"
	}`
	body, code := sendRequest(
		http.MethodPost,
		"/books",
		strings.NewReader(testBody),
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
	start(t)
	defer finish(t)

	// prepare data
	testBodyList := []string{
		`{
			"title": "TestTitleA",
			"author": "TestAuthorA",
			"publisher": "TestPublisherA"
		}`,
		`{
			"title": "TestTitleB",
			"author": "TestAuthorB",
			"publisher": "TestPublisherB"
		}`,
	}
	for _, testBody := range testBodyList {
		_, code := sendRequest(
			http.MethodPost,
			"/books",
			strings.NewReader(testBody),
			setupRouter(http.MethodPost, "/books", CreateBook),
		)
		assert.Equal(t, http.StatusOK, code)
	}

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
	assert.Equal(t, 2, len(res))
}

func TestGetBookByID(t *testing.T) {
	start(t)
	defer finish(t)

	// prepare data
	testBody := `{
		"title": "TestTitleA",
		"author": "TestAuthorA",
		"publisher": "TestPublisherA"
	}`
	_, code := sendRequest(
		http.MethodPost,
		"/books",
		strings.NewReader(testBody),
		setupRouter(http.MethodPost, "/books", CreateBook),
	)
	assert.Equal(t, http.StatusOK, code)

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
	start(t)
	defer finish(t)

	// prepare data
	testBody := `{
		"title": "TestTitleA",
		"author": "TestAuthorA",
		"publisher": "TestPublisherA"
	}`
	_, code := sendRequest(
		http.MethodPost,
		"/books",
		strings.NewReader(testBody),
		setupRouter(http.MethodPost, "/books", CreateBook),
	)
	assert.Equal(t, http.StatusOK, code)

	// test
	updateBody := `{
		"title": "TestTitleB",
		"author": "TestAuthorB",
		"publisher": "TestPublisherB"
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
	assert.Equal(t, "TestTitleB", res["title"].(string))
	assert.Equal(t, "TestAuthorB", res["author"].(string))
	assert.Equal(t, "TestPublisherB", res["publisher"].(string))
}

func TestDeleteBook(t *testing.T) {
	start(t)
	defer finish(t)

	// prepare data
	testBody := `{
		"title": "TestTitleA",
		"author": "TestAuthorA",
		"publisher": "TestPublisherA"
	}`
	_, code := sendRequest(
		http.MethodPost,
		"/books",
		strings.NewReader(testBody),
		setupRouter(http.MethodPost, "/books", CreateBook),
	)
	assert.Equal(t, http.StatusOK, code)

	// test
	_, code = sendRequest(
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

// start init test database
func start(t *testing.T) {
	assert.NoError(t, config.LoadConfig())
	assert.NoError(t, config.InitTestDB())
}

// start free test database
func finish(t *testing.T) {
	assert.NoError(t, config.FreeTestDB())
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
