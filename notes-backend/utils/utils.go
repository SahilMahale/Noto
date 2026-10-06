// Package utils holds generic, non-auth helper functions shared across the
// backend. Auth/JWT-specific helpers (key loading, token signing, the JWT
// middleware) and server middleware setup stay in server/utils.go — they're
// tightly coupled to that package's private key state and *notesService.
package utils

import (
	"os"
	"strings"
)

// GetEnvDefault returns the value of the given env var, or fallback if unset/empty.
func GetEnvDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// SystemTimeZoneName resolves the host's IANA timezone name (e.g. "Asia/Kolkata").
// Falls back to "UTC" when it can't be determined (e.g. no TZ env var and
// /etc/localtime isn't a zoneinfo symlink, common in minimal containers).
//
// Note: time.Local.String() can't be used for this — when Go derives the
// local zone from /etc/localtime directly (no TZ env var set), it hardcodes
// the Location's name to the literal string "Local", not the real zone name.
func SystemTimeZoneName() string {
	if tz := os.Getenv("TZ"); tz != "" {
		return tz
	}
	if target, err := os.Readlink("/etc/localtime"); err == nil {
		if i := strings.Index(target, "zoneinfo/"); i != -1 {
			return target[i+len("zoneinfo/"):]
		}
	}
	return "UTC"
}
