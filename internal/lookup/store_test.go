package lookup

import "testing"

func TestMemoryStoreContainsKnownURL(t *testing.T) {
	store := NewMemoryStore([]string{"malware.test/bad"})

	if !store.Contains("malware.test/bad") {
		t.Fatal("expected store to contain known URL")
	}
}

func TestMemoryStoreDoesNotContainUnknownURL(t *testing.T) {
	store := NewMemoryStore([]string{"malware.test/bad"})

	if store.Contains("example.com/safe") {
		t.Fatal("expected store not to contain unknown URL")
	}
}

func TestServiceReturnsMaliciousForKnownURL(t *testing.T) {
	service := NewService(NewMemoryStore([]string{"malware.test/bad"}))

	got := service.Lookup("malware.test/bad")

	if got.Verdict != VerdictMalicious {
		t.Fatalf("expected verdict %q, got %q", VerdictMalicious, got.Verdict)
	}
	if !got.Matched {
		t.Fatal("expected matched to be true")
	}
	if got.Reason == "" {
		t.Fatal("expected reason for malicious verdict")
	}
}

func TestServiceReturnsSafeForUnknownURL(t *testing.T) {
	service := NewService(NewMemoryStore([]string{"malware.test/bad"}))

	got := service.Lookup("example.com/safe")

	if got.Verdict != VerdictSafe {
		t.Fatalf("expected verdict %q, got %q", VerdictSafe, got.Verdict)
	}
	if got.Matched {
		t.Fatal("expected matched to be false")
	}
	if got.Reason != "" {
		t.Fatalf("expected empty reason, got %q", got.Reason)
	}
}
