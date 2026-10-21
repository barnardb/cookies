package main

import (
	"bytes"
	"encoding/json"
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

func TestStrongestCookie(t *testing.T) {
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
			if got := strongestCookie(tt.cookies).Value; got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWriteNetscape(t *testing.T) {
	expires := time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC)
	session := testCookie("session", "s", "example.com", "/", time.Time{})
	domain := testCookie("domain", "d", ".example.com", "/app", time.Time{})
	domain.Expires = expires
	secure := testCookie("secure", "x", "example.com", "/", time.Time{})
	secure.Secure = true
	httpOnly := testCookie("httponly", "h", ".example.com", "/", time.Time{})
	httpOnly.HttpOnly = true
	noPath := testCookie("nopath", "n", "example.com", "", time.Time{})

	var out bytes.Buffer
	writeNetscape(&out, []*kooky.Cookie{session, domain, secure, httpOnly, noPath})
	want := "# Netscape HTTP Cookie File\n\n" +
		"example.com\tFALSE\t/\tFALSE\t0\tsession\ts\n" + // session cookies have expiry 0, not a negative timestamp
		".example.com\tTRUE\t/app\tFALSE\t1893553445\tdomain\td\n" +
		"example.com\tFALSE\t/\tTRUE\t0\tsecure\tx\n" +
		"#HttpOnly_.example.com\tTRUE\t/\tFALSE\t0\thttponly\th\n" +
		"example.com\tFALSE\t/\tFALSE\t0\tnopath\tn\n"
	if got := out.String(); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestWriteJSON(t *testing.T) {
	c := testCookie("sid", "a<b>&c", "example.com", "/", time.Time{})

	var out bytes.Buffer
	if err := writeJSON(&out, []*kooky.Cookie{c}); err != nil {
		t.Fatal(err)
	}
	var array []map[string]any
	if err := json.Unmarshal(out.Bytes(), &array); err != nil {
		t.Fatalf("not a JSON array: %v\n%s", err, out.String())
	}
	if len(array) != 1 || array[0]["name"] != "sid" || array[0]["value"] != "a<b>&c" || array[0]["domain"] != "example.com" {
		t.Errorf("unexpected JSON: %s", out.String())
	}

	out.Reset()
	if err := writeJSON(&out, c); err != nil {
		t.Fatal(err)
	}
	var object map[string]any
	if err := json.Unmarshal(out.Bytes(), &object); err != nil {
		t.Fatalf("not a JSON object: %v\n%s", err, out.String())
	}
	if object["value"] != "a<b>&c" {
		t.Errorf("unexpected JSON: %s", out.String())
	}
}
