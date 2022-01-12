package postgres

import (
	"context"
	_ "embed" // embed schema
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Include the schema through embedding
//go:embed schema.sql
var schema string

// CurrentSchemaVersion is the current version of the schema
const CurrentSchemaVersion = 1

var (
// ErrNotInitialized is returned when the database
// schema is not migrated yet
)

// Status is the database / store status
type Status struct {
	Migrated        bool      `json:"migrated"`
	SchemaVersion   int       `json:"schema_version"`
	SchemaAppliedAt time.Time `json:"schema_applied_at"`
	Error           error     `json:"error"`
}

// Log writes the status into the log
func (s *Status) Log() {
	log.Println(
		"Database migrated:", s.Migrated)
	log.Println(
		"Schema version:", s.SchemaVersion,
		"applied at:", s.SchemaAppliedAt)
}

// The Manager supervises the database. It can migrate the
// schema and retrieve a status.
type Manager struct {
	pool *pgxpool.Pool
}

// NewManager creates a new database manager
func NewManager(pool *pgxpool.Pool) *Manager {
	return &Manager{
		pool: pool,
	}
}

// Start the background jobs for database management
func (m *Manager) Start(ctx context.Context) {
	m.Status(ctx).Log()
}

// Status retrieves the current schema version
// and checks if migrated. In case an error occures,
// it will be included in the result.
func (m *Manager) Status(ctx context.Context) *Status {
	status := &Status{}
	qry := `
		SELECT version, applied_at FROM __meta__
		 ORDER BY version DESC
		 LIMIT 1
	`
	err := m.pool.QueryRow(ctx, qry).Scan(
		&status.SchemaVersion,
		&status.SchemaAppliedAt,
	)
	if err != nil {
		status.Error = err
		return status
	}

	// Is the database migrated?
	status.Migrated = CurrentSchemaVersion == status.SchemaVersion

	return status
}

// Migrate applies the database intialisation script if required.
func (m *Manager) Migrate(ctx context.Context) error {
	s := m.Status(ctx)
	if s.Migrated {
		return nil
	}

	return m.Initialize(ctx)
}

// Initialize will apply the database schema. This will clear the
// database. However for now we treat the state as disposable.
func (m *Manager) Initialize(ctx context.Context) error {
	_, err := m.pool.Exec(ctx, schema)
	if err != nil {
		return err
	}
	return nil
}
