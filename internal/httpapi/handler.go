package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/mvillla/url-safety-checker/internal/lookup"
)

const lookupPrefix = "/urlinfo/1/"

type Handler struct {
	lookup *lookup.Service
}

func NewHandler(lookupService *lookup.Service) *Handler {
	return &Handler{lookup: lookupService}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.health)
	mux.HandleFunc("/readyz", h.ready)
	mux.HandleFunc(lookupPrefix, h.urlInfo)
	return mux
}

func (h *Handler) urlInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	hostport, path, ok := parseLookupPath(r.URL.EscapedPath())
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	normalizedURL, err := lookup.NormalizeKey(hostport, path, r.URL.RawQuery)
	if err != nil {
		if errors.Is(err, lookup.ErrInvalidURLKey) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, h.lookup.Lookup(normalizedURL))
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) ready(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func parseLookupPath(path string) (hostport, originalPath string, ok bool) {
	rest, ok := strings.CutPrefix(path, lookupPrefix)
	if !ok || rest == "" {
		return "", "", false
	}

	hostport, originalPath, found := strings.Cut(rest, "/")
	if !found || hostport == "" {
		return "", "", false
	}

	return hostport, "/" + originalPath, true
}
