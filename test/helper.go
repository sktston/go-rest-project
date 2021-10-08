package test

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rs/zerolog"
	"github.com/sktston/go-rest-project/database"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io"
	"net"
	"net/http/httptest"
	"testing"
	"time"
)

const postgresVersion = "13"
var testDBHost = ""
var testDBPort = ""

// CreatePostgres create postgres docker container
func CreatePostgres() (*dockertest.Pool, *dockertest.Resource, error) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, err
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        postgresVersion,
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, nil, err
	}
	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	_ = resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err := pool.Retry(func() error {
		testSqlDB, err := sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return testSqlDB.Ping()
	}); err != nil {
		return nil, nil, err
	}

	testDBHost, testDBPort, _ = net.SplitHostPort(hostAndPort)
	return pool, resource, nil
}

// RemovePostgres remove postgres docker container
func RemovePostgres(pool *dockertest.Pool, resource *dockertest.Resource) error {
	if err := pool.Purge(resource); err != nil {
		return err
	}
	return nil
}

// InitTestDB init test database
func InitTestDB(t *testing.T) {
	// prepare
	viper.Set("database.host", testDBHost)
	viper.Set("database.user", "user_name")
	viper.Set("database.password", "secret")
	viper.Set("database.dbname", "dbname")
	viper.Set("database.port", testDBPort)

	viper.Set("log.level", "TEST")

	assert.NoError(t, database.InitDB())
}

// FreeTestDB free test database
func FreeTestDB(t *testing.T) {
	assert.NoError(t, database.DropSchema(database.GetDB()))
}

// SetupRouter get router on given handler
func SetupRouter(httpMethod string, path string, handler gin.HandlerFunc) *gin.Engine {
	// prepare router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Handle(httpMethod, path, handler)
	return router
}

// SendRequest reads response from given http request.
func SendRequest(httpMethod string, target string, requestBody io.Reader, router *gin.Engine) (*bytes.Buffer, int) {
	// serve http on given response and request
	req := httptest.NewRequest(httpMethod, target, requestBody)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr.Body, rr.Code
}

// DisableLogger disable zero log and gin logger
func DisableLogger() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	gin.SetMode(gin.TestMode)
}