package config

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"net"
	"os"
	"testing"
	"time"
)

const (
	testDBUser     = "user_name"
	testDBPassword = "secret"
	testDBName     = "dbname"
)

var (
	testDBHost = ""
	testDBPort = ""
)

// Tests

func TestConnectingDatabase(t *testing.T) {
	// Open test DB with random prefix
	testDBPrefix := uuid.New().String()+"_"
	testDsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		testDBHost,
		testDBUser,
		testDBPassword,
		testDBName,
		testDBPort,
	)
	testDB, err := gorm.Open(postgres.Open(testDsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: testDBPrefix, // prefix is testId_
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	// Migrate the schema
	assert.NoError(t, MigrateSchema(testDB))

	// Test
	sqlDB, err := testDB.DB()
	assert.NoError(t, err)
	assert.NoError(t, sqlDB.Ping())

	// Free
	assert.NoError(t, DropSchema(testDB))
}

// Helpers

func TestMain(m *testing.M) {
	// create postgres docker
	pool, resource, err := createPostgres()
	if err != nil {
		log.Fatal().Msgf("Could not create postgres: %s", err)
	}

	//Run tests
	code := m.Run()

	// delete postgres docker
	if err := deletePostgres(pool, resource); err != nil {
		log.Fatal().Msgf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createPostgres() (*dockertest.Pool, *dockertest.Resource, error) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, err
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD="+ testDBPassword,
			"POSTGRES_USER="+ testDBUser,
			"POSTGRES_DB="+ testDBName,
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

// deletePostgres delete postgres docker
func deletePostgres(pool *dockertest.Pool, resource *dockertest.Resource) error {
	if err := pool.Purge(resource); err != nil {
		return err
	}
	return nil
}