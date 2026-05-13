package lookup

import (
	"errors"
	"net"
	"strings"
)

var ErrInvalidURLKey = errors.New("invalid URL key")

func NormalizeKey(hostport, path, rawQuery string) (string, error) {
	if hostport == "" || path == "" || !strings.HasPrefix(path, "/") {
		return "", ErrInvalidURLKey
	}

	normalizedHostport := normalizeHostport(hostport)
	key := normalizedHostport + path
	if rawQuery != "" {
		key += "?" + rawQuery
	}

	return key, nil
}

func normalizeHostport(hostport string) string {
	host, port, err := net.SplitHostPort(hostport)
	if err == nil {
		return strings.ToLower(host) + ":" + port
	}

	return strings.ToLower(hostport)
}
