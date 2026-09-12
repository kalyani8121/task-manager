package main

import (
	"database/sql"
	"testing"
)

func TestPostgresDriverRegistration(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost user=postgres dbname=taskmanager sslmode=disable")
	if err != nil {
		t.Fatalf("expected postgres driver to be registered, got %v", err)
	}
	defer db.Close()
}
