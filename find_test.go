package main

import (
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/browserutils/kooky"
)

func TestProfileMatcher(t *testing.T) {
	tests := []struct {
		name    string
		matcher profileMatcher
		profile string
		want    bool
	}{
		{"no selection matches everything", profileMatcher{}, "Default", true},
		{"exact name", profileMatcher{names: []string{"Profile 1"}}, "Profile 1", true},
		{"exact name is not a prefix match", profileMatcher{names: []string{"Profile 1"}}, "Profile 10", false},
		{"exact name is not a regex", profileMatcher{names: []string{"Profile ."}}, "Profile 1", false},
		{"any of several names", profileMatcher{names: []string{"Work", "Default"}}, "Default", true},
		{"pattern is unanchored", profileMatcher{patterns: []*regexp.Regexp{regexp.MustCompile("ork")}}, "Work", true},
		{"anchored pattern", profileMatcher{patterns: []*regexp.Regexp{regexp.MustCompile("^Profile [0-9]$")}}, "Profile 10", false},
		{"name or pattern", profileMatcher{names: []string{"Work"}, patterns: []*regexp.Regexp{regexp.MustCompile("^Profile")}}, "Profile 2", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.matcher.matches(tt.profile); got != tt.want {
				t.Errorf("matches(%q) = %v, want %v", tt.profile, got, tt.want)
			}
		})
	}
}

func TestAppliesToURLAndNameAtTime(t *testing.T) {
	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
	cookie := func(domain, path string, edit ...func(*kooky.Cookie)) *kooky.Cookie {
		c := &kooky.Cookie{}
		c.Name, c.Domain, c.Path = "sid", domain, path
		for _, e := range edit {
			e(c)
		}
		return c
	}
	secure := func(c *kooky.Cookie) { c.Secure = true }
	expires := func(t time.Time) func(*kooky.Cookie) { return func(c *kooky.Cookie) { c.Expires = t } }

	tests := []struct {
		name   string
		url    string
		cookie *kooky.Cookie
		want   bool
	}{
		// domain matching (RFC 6265 §5.1.3)
		{"host-only cookie, same host", "http://example.com/", cookie("example.com", "/"), true},
		{"host-only cookie, subdomain", "http://www.example.com/", cookie("example.com", "/"), false},
		{"domain cookie, same host", "http://example.com/", cookie(".example.com", "/"), true},
		{"domain cookie, subdomain", "http://a.b.example.com/", cookie(".example.com", "/"), true},
		{"domain cookie, suffix that isn't a subdomain", "http://notexample.com/", cookie(".example.com", "/"), false},
		{"different host", "http://example.org/", cookie("example.com", "/"), false},
		{"URL with port", "http://localhost:8080/", cookie("localhost", "/"), true},
		{"URL with port, domain cookie", "https://www.example.com:8443/", cookie(".example.com", "/"), true},
		{"URL host is case-insensitive", "http://WWW.Example.COM/", cookie(".example.com", "/"), true},

		// path matching (RFC 6265 §5.1.4)
		{"root cookie path", "http://example.com/a/b", cookie("example.com", "/"), true},
		{"identical path", "http://example.com/foo", cookie("example.com", "/foo"), true},
		{"cookie path is a directory prefix", "http://example.com/foo/bar", cookie("example.com", "/foo"), true},
		{"cookie path ending in slash", "http://example.com/foo/bar", cookie("example.com", "/foo/"), true},
		{"cookie path is a string prefix but not a directory prefix", "http://example.com/foobar", cookie("example.com", "/foo"), false},
		{"request path shorter than cookie path", "http://example.com/fo", cookie("example.com", "/foo"), false},
		{"empty URL path is treated as /", "http://example.com", cookie("example.com", "/foo"), false},
		{"empty URL path matches root cookie", "http://example.com", cookie("example.com", "/"), true},
		{"empty cookie path", "http://example.com/foo", cookie("example.com", ""), true},
		{"query string isn't part of the path", "http://example.com/foo?x=/bar", cookie("example.com", "/foo"), true},

		// secure and expiry
		{"secure cookie over https", "https://example.com/", cookie("example.com", "/", secure), true},
		{"secure cookie over http", "http://example.com/", cookie("example.com", "/", secure), false},
		{"session cookie (no expiry)", "http://example.com/", cookie("example.com", "/"), true},
		{"unexpired cookie", "http://example.com/", cookie("example.com", "/", expires(now.Add(time.Hour))), true},
		{"expired cookie", "http://example.com/", cookie("example.com", "/", expires(now.Add(-time.Hour))), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.url)
			if err != nil {
				t.Fatal(err)
			}
			if got := appliesToURLAndNameAtTime(u, "", now, nil).Filter(tt.cookie); got != tt.want {
				t.Errorf("cookie %s%s for %s: got %v, want %v", tt.cookie.Domain, tt.cookie.Path, tt.url, got, tt.want)
			}
		})
	}
}

func TestAppliesToURLAndNameAtTimeFiltersByName(t *testing.T) {
	u, _ := url.Parse("http://example.com/")
	c := &kooky.Cookie{}
	c.Name, c.Domain, c.Path = "sid", "example.com", "/"
	if !appliesToURLAndNameAtTime(u, "sid", time.Now(), nil).Filter(c) {
		t.Error("cookie with the requested name was rejected")
	}
	if appliesToURLAndNameAtTime(u, "other", time.Now(), nil).Filter(c) {
		t.Error("cookie with a different name was accepted")
	}
}
