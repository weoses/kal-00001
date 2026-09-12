package functional_test

import (
	"context"
	"testing"
)

// TestSearchMeme_MultipleAccountIds_ReturnsBoth verifies that SearchMeme with
// two account ids returns memes from either account, while a single-account
// search only sees its own account's memes.
func TestSearchMeme_MultipleAccountIds_ReturnsBoth(t *testing.T) {
	resetState(t)
	ctx := context.Background()
	c := testClient()
	acctA := accountID("da001")
	acctB := accountID("da002")

	if _, err := c.CreateMemeFromBytes(ctx, acctA, minimalJPEG(1)); err != nil {
		t.Fatalf("create meme for acctA: %v", err)
	}
	if _, err := c.CreateMemeFromBytes(ctx, acctB, minimalJPEG(2)); err != nil {
		t.Fatalf("create meme for acctB: %v", err)
	}

	respA, err := c.SearchMemePage(ctx, acctA, "", 10, nil)
	if err != nil {
		t.Fatalf("SearchMeme acctA only: %v", err)
	}
	if len(respA.GetResults()) != 1 {
		t.Fatalf("expected 1 result for acctA alone, got %d", len(respA.GetResults()))
	}

	respBoth, err := c.SearchMemePageMulti(ctx, []string{acctA, acctB}, "", 10, nil)
	if err != nil {
		t.Fatalf("SearchMeme both accounts: %v", err)
	}
	if len(respBoth.GetResults()) != 2 {
		t.Fatalf("expected 2 results across both accounts, got %d", len(respBoth.GetResults()))
	}
}

// TestSearchMeme_MultipleAccountIds_IdLookupMatchesEitherAccount verifies
// that the exact-UUID search branch (IdSearcher) also honors every provided
// account id, not just the first.
func TestSearchMeme_MultipleAccountIds_IdLookupMatchesEitherAccount(t *testing.T) {
	resetState(t)
	ctx := context.Background()
	c := testClient()
	acctA := accountID("da003")
	acctB := accountID("da004")

	created, err := c.CreateMemeFromBytes(ctx, acctB, minimalJPEG(3))
	if err != nil {
		t.Fatalf("create meme for acctB: %v", err)
	}
	memeID := created.GetResult().GetId()
	if memeID == "" {
		t.Fatal("expected non-empty created meme id")
	}

	// Searching by the meme's own id, authorized for acctA+acctB, should find
	// it even though it actually belongs to acctB.
	resp, err := c.SearchMemePageMulti(ctx, []string{acctA, acctB}, memeID, 10, nil)
	if err != nil {
		t.Fatalf("SearchMeme by id across both accounts: %v", err)
	}
	if len(resp.GetResults()) != 1 || resp.GetResults()[0].GetId() != memeID {
		t.Fatalf("expected id-lookup to find meme %s, got %+v", memeID, resp.GetResults())
	}

	// Searching by that id with only acctA (which doesn't own it) must not find it.
	respAOnly, err := c.SearchMemePage(ctx, acctA, memeID, 10, nil)
	if err != nil {
		t.Fatalf("SearchMeme by id acctA only: %v", err)
	}
	if len(respAOnly.GetResults()) != 0 {
		t.Fatalf("expected id-lookup restricted to acctA to find nothing, got %+v", respAOnly.GetResults())
	}
}

// TestGetRandomMeme_MultipleAccountIds_ReturnsFromEither verifies GetRandomMeme
// can return a meme belonging to any of the given account ids.
func TestGetRandomMeme_MultipleAccountIds_ReturnsFromEither(t *testing.T) {
	resetState(t)
	ctx := context.Background()
	c := testClient()
	acctA := accountID("da005")
	acctB := accountID("da006")

	if _, err := c.CreateMemeFromBytes(ctx, acctB, minimalJPEG(4)); err != nil {
		t.Fatalf("create meme for acctB: %v", err)
	}

	// acctA alone has nothing.
	_, err := c.GetRandomMemeMulti(ctx, []string{acctA}, "")
	if err == nil {
		t.Fatal("expected an error/not-found for acctA alone with no memes")
	}

	resp, err := c.GetRandomMemeMulti(ctx, []string{acctA, acctB}, "")
	if err != nil {
		t.Fatalf("GetRandomMeme across both accounts: %v", err)
	}
	if resp.GetResult().GetId() == "" {
		t.Fatal("expected a meme from acctB via the combined account id list")
	}
}
