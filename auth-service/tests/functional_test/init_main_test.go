package functional_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	tcPostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/weoses/memelo/auth-service/conf"
	commonconfig "github.com/weoses/memelo/common/config"
)

const (
	testDbUser     = "auth_test"
	testDbPassword = "auth_test"
	testDbName     = "auth_test"
)

var (
	testAddr  string
	testStop  func()
	testSqlDb *sql.DB
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := tcPostgres.Run(ctx, "postgres:16-alpine",
		tcPostgres.WithDatabase(testDbName),
		tcPostgres.WithUsername(testDbUser),
		tcPostgres.WithPassword(testDbPassword),
		tcPostgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatalf("start postgres: %v", err)
	}
	defer func() { _ = pgContainer.Terminate(ctx) }()

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("postgres connection string: %v", err)
	}

	testSqlDb, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("open sql db: %v", err)
	}
	defer func() { _ = testSqlDb.Close() }()

	cfg := buildTestConfig(dsn)

	// Starting the app runs auth-service's embedded migrations (app.Module()
	// invokes storage.RunMigrations at startup), so the schema is ready by
	// the time startTestServer returns.
	testAddr, testStop, err = startTestServer(cfg)
	if err != nil {
		log.Fatalf("start test server: %v", err)
	}
	defer testStop()

	os.Exit(m.Run())
}

func buildTestConfig(dsn string) *conf.Config {
	return &conf.Config{
		Server:   &commonconfig.ServerConfig{ListenAddress: ":0"},
		Log:      &commonconfig.LoggingConfig{Level: "warn"},
		Postgres: &conf.PostgresConfig{Dsn: dsn},
	}
}

// resetState clears all per-test data between tests, keeping the
// migration-seeded permissions table intact.
func resetState(t *testing.T) {
	t.Helper()
	if _, err := testSqlDb.ExecContext(context.Background(), "TRUNCATE users CASCADE"); err != nil {
		t.Fatalf("reset state: %v", err)
	}
}

// testClient returns a ready-to-use TestClient pointed at the running test server.
func testClient() *TestClient {
	return NewTestClient(testAddr)
}
