//go:build e2e
// +build e2e

package e2e

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"

	_ "github.com/lib/pq"
)

func TestConnectionToServer(t *testing.T) {
	id := rand.Int()

	databaseName := fmt.Sprintf("test_db_%d", id)
	dbUserName := fmt.Sprintf("test_user_%d", id)
	dbUserPassword := fmt.Sprintf("test_user_password_%d", id)

	if err := provisionDbForUser(databaseName, dbUserName, dbUserPassword); err != nil {
		t.Fatalf("Failed to provision database: %v", err)
	}

	db, err := connectToDb("app", "54321", dbUserName, dbUserPassword, databaseName)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("SELECT 1")
	if err != nil {
		t.Fatalf("Failed to execute query: %v", err)
	}
}

func provisionDbForUser(dbName string, userName string, userPass string) error {
	db, err := connectToDb("postgres", "5432", "postgres", "postgres", "postgres")
	if err != nil {
		return err
	}
	defer db.Close()

	createDbQuery := fmt.Sprintf("CREATE DATABASE %s", dbName)
	if _, err := db.Exec(createDbQuery); err != nil {
		return err
	}

	createUserQuery := fmt.Sprintf("CREATE USER %s WITH ENCRYPTED PASSWORD '%s'", userName, userPass)
	if _, err := db.Exec(createUserQuery); err != nil {
		return err
	}

	grantPrivilegesQuery := fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", dbName, userName)
	if _, err := db.Exec(grantPrivilegesQuery); err != nil {
		return err
	}

	return nil
}

func connectToDb(host string, port string, user string, password string, dbname string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)
	return sql.Open("postgres", connStr)
}
