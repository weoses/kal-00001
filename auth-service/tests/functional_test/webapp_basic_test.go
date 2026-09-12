package functional_test

import (
	"context"
	"testing"

	v1 "github.com/weoses/memelo/gen/proto/v1"
)

func TestWebappBasicAuthorize_CorrectCredentials_ReturnsOkWithPermissions(t *testing.T) {
	resetState(t)
	ctx := context.Background()
	c := testClient()

	userId := seedUser(t, "PUBLISH")
	linkWebappBasic(t, userId, "alice", "correct-horse-battery-staple")

	resp, err := c.AuthorizeWebappBasic(ctx, "alice", "correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("authorize webapp basic: %v", err)
	}

	result := resp.GetResult()
	if result.GetStatus() != v1.AuthStatus_AUTH_STATUS_OK {
		t.Fatalf("expected status OK, got %v", result.GetStatus())
	}
	if result.GetUserId() != userId {
		t.Fatalf("expected user id %q, got %q", userId, result.GetUserId())
	}
	assertContainsAll(t, result.GetPermissions(), "PUBLISH")
}

func TestWebappBasicAuthorize_WrongPassword_ReturnsNotFound(t *testing.T) {
	resetState(t)
	ctx := context.Background()
	c := testClient()

	userId := seedUser(t, "PUBLISH")
	linkWebappBasic(t, userId, "bob", "correct-password")

	resp, err := c.AuthorizeWebappBasic(ctx, "bob", "wrong-password")
	if err != nil {
		t.Fatalf("authorize webapp basic: %v", err)
	}

	result := resp.GetResult()
	if result.GetStatus() != v1.AuthStatus_AUTH_STATUS_NOT_FOUND {
		t.Fatalf("expected status NOT_FOUND for a wrong password, got %v", result.GetStatus())
	}
	if len(result.GetPermissions()) != 0 {
		t.Fatalf("expected no permissions for a rejected login, got %v", result.GetPermissions())
	}
}

func TestWebappBasicAuthorize_UnknownUsername_ReturnsNotFound(t *testing.T) {
	resetState(t)
	ctx := context.Background()
	c := testClient()

	resp, err := c.AuthorizeWebappBasic(ctx, "nobody", "irrelevant")
	if err != nil {
		t.Fatalf("authorize webapp basic: %v", err)
	}

	result := resp.GetResult()
	if result.GetStatus() != v1.AuthStatus_AUTH_STATUS_NOT_FOUND {
		t.Fatalf("expected status NOT_FOUND for an unknown username, got %v", result.GetStatus())
	}
}
