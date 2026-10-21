package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/browserutils/kooky"
)

func testCookie(name, value, domain, path string, created time.Time) *kooky.Cookie {
	c := &kooky.Cookie{Creation: created}
	c.Name, c.Value, c.Domain, c.Path = name, value, domain, path
	return c
}

func TestFormatCookies(t *testing.T) {
	var out bytes.Buffer
	formatCookies(&out, []*kooky.Cookie{
		testCookie("a", "1", "example.com", "/", time.Time{}),
		testCookie("b", "2=3", "example.com", "/", time.Time{}),
	})
	if got, want := out.String(), "a=1;b=2=3"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWriteStrongestValue(t *testing.T) {
	older := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	newer := older.Add(time.Hour)
	tests := []struct {
		name    string
		cookies []*kooky.Cookie
		want    string
	}{
		{"single cookie", []*kooky.Cookie{
			testCookie("sid", "only", "example.com", "/", older),
		}, "only"},
		{"longer domain wins", []*kooky.Cookie{
			testCookie("sid", "parent", ".example.com", "/", newer),
			testCookie("sid", "sub", ".www.example.com", "/", older),
		}, "sub"},
		{"longer path wins for same domain", []*kooky.Cookie{
			testCookie("sid", "root", "example.com", "/", newer),
			testCookie("sid", "deep", "example.com", "/foo", older),
		}, "deep"},
		{"host-only beats domain cookie for the same domain", []*kooky.Cookie{
			testCookie("sid", "domain", ".example.com", "/", newer),
			testCookie("sid", "host-only", "example.com", "/", older),
		}, "host-only"},
		{"host-only beats domain cookie regardless of order", []*kooky.Cookie{
			testCookie("sid", "host-only", "example.com", "/", older),
			testCookie("sid", "domain", ".example.com", "/foo", newer),
		}, "host-only"},
		{"longer domain still beats host-only", []*kooky.Cookie{
			testCookie("sid", "host-only", "example.com", "/", older),
			testCookie("sid", "sub", ".www.example.com", "/", older),
		}, "sub"},
		{"newer wins for same domain and path", []*kooky.Cookie{
			testCookie("sid", "old", "example.com", "/", older),
			testCookie("sid", "new", "example.com", "/", newer),
		}, "new"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			writeStrongestValue(&out, tt.cookies)
			if got := out.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
