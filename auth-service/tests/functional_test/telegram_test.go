package functional_test

import (
	"context"
	"testing"

	v1 "github.com/weoses/memelo/gen/proto/v1"
)

func TestTelegramAuthorize_LinkedUser_ReturnsOkWithPermissions(t *testing.T) {
	resetState(t)
	ctx := context.Background()
	c := testClient()

	userId := seedUser(t, "CREATE", "DELETE")
	linkTelegram(t, userId, 111)

	resp, err := c.AuthorizeTelegram(ctx, 111)
	if err != nil {
		t.Fatalf("authorize telegram: %v", err)
	}

	result := resp.GetResult()
	if result.GetStatus() != v1.AuthStatus_AUTH_STATUS_OK {
		t.Fatalf("expected status OK, got %v", result.GetStatus())
	}
	if result.GetUserId() != userId {
		t.Fatalf("expected user id %q, got %q", userId, result.GetUserId())
	}
	if result.GetAccountId() == "" {
		t.Fatal("expected non-empty account id")
	}
	assertContainsAll(t, result.GetPermissions(), "CREATE", "DELETE")
}

func TestTelegramAuthorize_UnlinkedId_ReturnsNotFound(t *testing.T) {
	resetState(t)
	ctx := context.Background()
	c := testClient()

	// A user exists, but is linked to a different Telegram id.
	userId := seedUser(t, "CREATE")
	linkTelegram(t, userId, 222)

	resp, err := c.AuthorizeTelegram(ctx, 999)
	if err != nil {
		t.Fatalf("authorize telegram: %v", err)
	}

	result := resp.GetResult()
	if result.GetStatus() != v1.AuthStatus_AUTH_STATUS_NOT_FOUND {
		t.Fatalf("expected status NOT_FOUND, got %v", result.GetStatus())
	}
	if len(result.GetPermissions()) != 0 {
		t.Fatalf("expected no permissions for an unlinked id, got %v", result.GetPermissions())
	}
}

func TestTelegramAuthorize_SecondLinkedId_TreatedAsSameUser(t *testing.T) {
	resetState(t)
	ctx := context.Background()
	c := testClient()

	// One user, two personal Telegram accounts.
	userId := seedUser(t, "PUBLISH")
	linkTelegram(t, userId, 333, 444)

	for _, telegramId := range []int64{333, 444} {
		resp, err := c.AuthorizeTelegram(ctx, telegramId)
		if err != nil {
			t.Fatalf("authorize telegram id %d: %v", telegramId, err)
		}
		result := resp.GetResult()
		if result.GetStatus() != v1.AuthStatus_AUTH_STATUS_OK {
			t.Fatalf("telegram id %d: expected status OK, got %v", telegramId, result.GetStatus())
		}
		if result.GetUserId() != userId {
			t.Fatalf("telegram id %d: expected user id %q, got %q", telegramId, userId, result.GetUserId())
		}
	}
}

func assertContainsAll(t *testing.T, actual []string, expected ...string) {
	t.Helper()
	set := make(map[string]struct{}, len(actual))
	for _, v := range actual {
		set[v] = struct{}{}
	}
	for _, want := range expected {
		if _, ok := set[want]; !ok {
			t.Fatalf("expected %q in %v", want, actual)
		}
	}
}
