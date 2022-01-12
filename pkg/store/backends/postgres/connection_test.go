package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

// ConnectTest uses defaults for the connection
// environment variable is not set.
func ConnectTest() *pgxpool.Pool {
	ctx := context.Background()
	url := os.Getenv("ALICE_TEST_DB_URL")
	if url == "" {
		url = "postgres://postgres:postgres@localhost:5432/alice_test"
	}
	p, err := Connect(ctx, &ConnectOpts{
		URL:      url,
		MinConns: 2,
		MaxConns: 16})
	if err != nil {
		panic(err)
	}

	m := NewManager(p)
	err = m.Initialize(ctx)
	if err != nil {
		panic(err)
	}

	return p
}

func TestInitialize(t *testing.T) {
	p := ConnectTest()
	m := NewManager(p)
	s := m.Status(context.Background())
	if s.Error != nil {
		t.Error(s.Error)
	}
	if s.Migrated == false {
		t.Error("schema is not migrated, current:", s.SchemaVersion)
	}
}
