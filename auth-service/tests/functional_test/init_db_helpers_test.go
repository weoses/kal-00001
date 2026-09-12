package functional_test

import (
	"context"
	"testing"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// seedUser inserts a user (letting Postgres generate its id/account_id) and
// grants it the given permission codes, returning the new user's id.
func seedUser(t *testing.T, permissionCodes ...string) string {
	t.Helper()
	ctx := context.Background()

	var userId string
	if err := testSqlDb.QueryRowContext(ctx, "INSERT INTO users DEFAULT VALUES RETURNING id").Scan(&userId); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	for _, code := range permissionCodes {
		if _, err := testSqlDb.ExecContext(ctx,
			"INSERT INTO user_permissions (user_id, permission_id) SELECT $1, id FROM permissions WHERE code = $2",
			userId, code,
		); err != nil {
			t.Fatalf("grant permission %s: %v", code, err)
		}
	}

	return userId
}

// linkTelegram links the given Telegram numeric ids to userId.
func linkTelegram(t *testing.T, userId string, telegramIds ...int64) {
	t.Helper()
	if _, err := testSqlDb.ExecContext(context.Background(),
		"INSERT INTO integration_telegram (user_id, telegram_ids) VALUES ($1, $2)",
		userId, pq.Int64Array(telegramIds),
	); err != nil {
		t.Fatalf("link telegram: %v", err)
	}
}

// linkWebappBasic links a username/password to userId, hashing the password
// the same way auth-service does.
func linkWebappBasic(t *testing.T, userId string, username string, password string) {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if _, err := testSqlDb.ExecContext(context.Background(),
		"INSERT INTO integration_webapp_basic (user_id, username, password_hash) VALUES ($1, $2, $3)",
		userId, username, string(hash),
	); err != nil {
		t.Fatalf("link webapp basic: %v", err)
	}
}
