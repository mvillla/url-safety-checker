package lookup

import (
	"errors"
	"testing"
)

func TestNormalizeKeyLowercasesHost(t *testing.T) {
	got, err := NormalizeKey("Example.COM", "/bad", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := "example.com/bad"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestNormalizeKeyPreservesPathCase(t *testing.T) {
	got, err := NormalizeKey("example.com", "/Path/To/Page", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := "example.com/Path/To/Page"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestNormalizeKeyPreservesQueryCase(t *testing.T) {
	got, err := NormalizeKey("example.com", "/path", "Token=ABC&x=1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := "example.com/path?Token=ABC&x=1"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestNormalizeKeyPreservesPort(t *testing.T) {
	got, err := NormalizeKey("Example.COM:443", "/path", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := "example.com:443/path"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestNormalizeKeyRejectsEmptyHost(t *testing.T) {
	_, err := NormalizeKey("", "/path", "")
	if !errors.Is(err, ErrInvalidURLKey) {
		t.Fatalf("expected ErrInvalidURLKey, got %v", err)
	}
}

func TestNormalizeKeyRejectsEmptyPath(t *testing.T) {
	_, err := NormalizeKey("example.com", "", "")
	if !errors.Is(err, ErrInvalidURLKey) {
		t.Fatalf("expected ErrInvalidURLKey, got %v", err)
	}
}

func TestNormalizeKeyRejectsPathWithoutSlash(t *testing.T) {
	_, err := NormalizeKey("example.com", "path", "")
	if !errors.Is(err, ErrInvalidURLKey) {
		t.Fatalf("expected ErrInvalidURLKey, got %v", err)
	}
}
