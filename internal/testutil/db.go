package testutil

import (
    "context"
    "fmt"
    "testing"

    "github.com/jackc/pgx/v5/pgxpool"
)

// TestDBConfig holds database configuration for tests
type TestDBConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

// DefaultTestDBConfig returns default test database configuration
func DefaultTestDBConfig() TestDBConfig {
    return TestDBConfig{
        Host:     "localhost",
        Port:     "5432",
        User:     "u0_a283",
        Password: "",
        DBName:   "mydb_test",
        SSLMode:  "disable",
    }
}

// ConnectionString returns the PostgreSQL connection string
func (c TestDBConfig) ConnectionString() string {
    return fmt.Sprintf(
        "postgres://%s@%s:%s/%s?sslmode=%s",
        c.User, c.Host, c.Port, c.DBName, c.SSLMode,
    )
}

// SetupTestDB creates a test database connection and cleans up
func SetupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
    config := DefaultTestDBConfig()
    
    // Create test database if it doesn't exist
    createTestDB(t, config)
    
    // Connect to test database
    conn, err := pgxpool.New(context.Background(), config.ConnectionString())
    if err != nil {
        t.Fatalf("Failed to connect to test database: %v", err)
    }
    
    // Cleanup function
    cleanup := func() {
        conn.Close()
        dropTestDB(t, config)
    }
    
    return conn, cleanup
}

// createTestDB creates the test database
func createTestDB(t *testing.T, config TestDBConfig) {
    // Connect to default postgres database
    connStr := fmt.Sprintf(
        "postgres://%s@%s:%s/postgres?sslmode=%s",
        config.User, config.Host, config.Port, config.SSLMode,
    )
    
    conn, err := pgxpool.New(context.Background(), connStr)
    if err != nil {
        t.Skipf("Skipping integration tests: %v", err)
        return
    }
    defer conn.Close()
    
    // Drop test database if exists
    _, err = conn.Exec(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %s", config.DBName))
    if err != nil {
        t.Logf("Warning: Could not drop test database: %v", err)
    }
    
    // Create test database
    _, err = conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", config.DBName))
    if err != nil {
        t.Fatalf("Failed to create test database: %v", err)
    }
}

// dropTestDB drops the test database
func dropTestDB(t *testing.T, config TestDBConfig) {
    connStr := fmt.Sprintf(
        "postgres://%s@%s:%s/postgres?sslmode=%s",
        config.User, config.Host, config.Port, config.SSLMode,
    )
    
    conn, err := pgxpool.New(context.Background(), connStr)
    if err != nil {
        t.Logf("Warning: Could not connect to drop test database: %v", err)
        return
    }
    defer conn.Close()
    
    _, err = conn.Exec(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %s", config.DBName))
    if err != nil {
        t.Logf("Warning: Could not drop test database: %v", err)
    }
}

// RunMigrations runs migrations on the test database
func RunMigrations(t *testing.T, conn *pgxpool.Pool) {
    // Create goose table
    _, err := conn.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS goose_db_version (
            id SERIAL PRIMARY KEY,
            version_id BIGINT NOT NULL,
            is_applied BOOLEAN NOT NULL,
            tstamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        t.Fatalf("Failed to create migrations table: %v", err)
    }
    
    // Run migrations
    migrations := []string{
        `CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            email TEXT UNIQUE NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            status VARCHAR(20) DEFAULT 'active',
            last_login TIMESTAMP
        )`,
        `CREATE TABLE IF NOT EXISTS posts (
            id SERIAL PRIMARY KEY,
            user_id INTEGER NOT NULL,
            title TEXT NOT NULL,
            content TEXT,
            published BOOLEAN DEFAULT FALSE,
            FOREIGN KEY (user_id) REFERENCES users(id)
        )`,
        `CREATE TABLE IF NOT EXISTS comments (
            id SERIAL PRIMARY KEY,
            post_id INTEGER NOT NULL,
            user_id INTEGER NOT NULL,
            content TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (post_id) REFERENCES posts(id),
            FOREIGN KEY (user_id) REFERENCES users(id)
        )`,
        `CREATE INDEX idx_comments_post_id ON comments(post_id)`,
        `CREATE INDEX idx_comments_user_id ON comments(user_id)`,
    }
    
    for _, migration := range migrations {
        _, err := conn.Exec(context.Background(), migration)
        if err != nil {
            t.Fatalf("Failed to run migration: %v", err)
        }
    }
}

// CleanDB cleans all tables in the test database
func CleanDB(t *testing.T, conn *pgxpool.Pool) {
    tables := []string{"comments", "posts", "users"}
    for _, table := range tables {
        _, err := conn.Exec(context.Background(), fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table))
        if err != nil {
            t.Logf("Warning: Could not truncate table %s: %v", table, err)
        }
    }
}
