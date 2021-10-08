package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"net"
	"os"
	"testing"
	"time"
)

const postgresVersion = "13"
var testDBHost = ""
var testDBPort = ""

// Tests

func TestInitDB(t *testing.T) {
	// prepare
	viper.Set("database.host", testDBHost)
	viper.Set("database.user", "user_name")
	viper.Set("database.password", "secret")
	viper.Set("database.dbname", "dbname")
	viper.Set("database.port", testDBPort)

	viper.Set("log.level", "TEST")

	// Test
	assert.NoError(t, InitDB())

	sqlDB, err := GetDB().DB()
	assert.NoError(t, err)
	assert.NoError(t, sqlDB.Ping())

	// Free
	assert.NoError(t, DropSchema(GetDB()))
}

// Helpers

// TestMain main function with postgres database
func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	// create postgres docker container
	pool, resource, err := createPostgres()
	if err != nil {
		fmt.Printf("Could not create postgres: %s", err)
		os.Exit(-1)
	}

	//Run tests
	code := m.Run()

	// remove postgres docker container
	if err := removePostgres(pool, resource); err != nil {
		fmt.Printf("Could not remove postgres: %s", err)
		os.Exit(-1)
	}

	os.Exit(code)
}

// createPostgres create postgres docker container
func createPostgres() (*dockertest.Pool, *dockertest.Resource, error) {
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

// removePostgres remove postgres docker container
func removePostgres(pool *dockertest.Pool, resource *dockertest.Resource) error {
	if err := pool.Purge(resource); err != nil {
		return err
	}
	return nil
}