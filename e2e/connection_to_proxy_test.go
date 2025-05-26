//go:build e2e
// +build e2e

package e2e

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestConnectionToServer(t *testing.T) {
	connStr := "host=pg port=5432 user=e2e_test_user password=e2e_test_password dbname=e2e_test_db sslmode=require"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	row, qErr := db.Query("SELECT 1")
	if qErr != nil {
		t.Fatalf("Failed to execute query: %v", qErr)
	}

	assert.Equal(t, row.Next(), true, "Expected a row to be returned")
	var result int
	if err := row.Scan(&result); err != nil {
		t.Fatalf("Failed to scan row: %v", err)
	}
	assert.Equal(t, result, 1, "Expected result to be 1")
}
