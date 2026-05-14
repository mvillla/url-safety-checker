package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mvillla/url-safety-checker/internal/lookup"
)

func TestHealthz(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestReadyz(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestHealthzRejectsUnsupportedMethod(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestURLInfoReturnsMaliciousForKnownURL(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/urlinfo/1/malware.test/bad", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got lookup.Result
	decodeJSON(t, rec, &got)

	if got.Verdict != lookup.VerdictMalicious {
		t.Fatalf("expected verdict %q, got %q", lookup.VerdictMalicious, got.Verdict)
	}
	if got.NormalizedURL != "malware.test/bad" {
		t.Fatalf("expected normalized URL %q, got %q", "malware.test/bad", got.NormalizedURL)
	}
}

func TestURLInfoReturnsSafeForUnknownURL(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/urlinfo/1/example.com/safe", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got lookup.Result
	decodeJSON(t, rec, &got)

	if got.Verdict != lookup.VerdictSafe {
		t.Fatalf("expected verdict %q, got %q", lookup.VerdictSafe, got.Verdict)
	}
	if got.Matched {
		t.Fatal("expected matched to be false")
	}
}

func TestURLInfoPreservesQueryString(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/urlinfo/1/malware.test/phishing?Token=ABC", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got lookup.Result
	decodeJSON(t, rec, &got)

	if got.NormalizedURL != "malware.test/phishing?Token=ABC" {
		t.Fatalf("expected normalized URL %q, got %q", "malware.test/phishing?Token=ABC", got.NormalizedURL)
	}
	if got.Verdict != lookup.VerdictMalicious {
		t.Fatalf("expected verdict %q, got %q", lookup.VerdictMalicious, got.Verdict)
	}
}

func TestURLInfoRejectsMalformedPath(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/urlinfo/1/example.com", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestURLInfoRejectsUnsupportedMethod(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/urlinfo/1/malware.test/bad", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func newTestHandler() http.Handler {
	store := lookup.NewMemoryStore([]string{
		"malware.test/bad",
		"malware.test/phishing?Token=ABC",
	})
	service := lookup.NewService(store)
	return NewHandler(service).Routes()
}

func decodeJSON(t *testing.T, rec *httptest.ResponseRecorder, v any) {
	t.Helper()

	if err := json.NewDecoder(rec.Body).Decode(v); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}
