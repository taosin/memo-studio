package database_test

import (
	"testing"

	"memo-studio/backend/database"
)

func TestInitCreatesAdminFromEnv(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("MEMO_DB_PATH", tmp+"/notes.db")
	t.Setenv("MEMO_ADMIN_PASSWORD", "StrongPass123!")

	if err := database.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}

	var cnt int
	if err := database.DB.QueryRow(`SELECT COUNT(1) FROM users WHERE username='admin' AND is_admin=1`).Scan(&cnt); err != nil {
		t.Fatalf("query admin: %v", err)
	}
	if cnt != 1 {
		t.Fatalf("expected admin user, got %d", cnt)
	}

	var must int
	if err := database.DB.QueryRow(`SELECT must_change_password FROM users WHERE username='admin'`).Scan(&must); err != nil {
		t.Fatalf("query must_change_password: %v", err)
	}
	if must != 1 {
		t.Fatalf("expected must_change_password=1, got %d", must)
	}
}

