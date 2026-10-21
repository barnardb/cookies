package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/browserutils/kooky"
)

const (
	formatHeader   = "header"
	formatNetscape = "netscape"
	formatJSON     = "json"
)

var formats = []string{formatHeader, formatNetscape, formatJSON}

func formatCookies(w io.Writer, cookies []*kooky.Cookie) {
	for i, cookie := range cookies {
		if i > 0 {
			fmt.Fprint(w, ";")
		}
		formatCookie(w, cookie)
	}
}

func formatCookie(w io.Writer, cookie *kooky.Cookie) {
	fmt.Fprintf(w, "%s=%s", cookie.Name, cookie.Value)
}

// writeNetscape writes cookies in the Netscape cookies.txt format read by curl (-b), wget
// (--load-cookies) and yt-dlp (--cookies). We don't use kooky.ExportCookies because it writes
// session cookies with a large negative expiry, which curl treats as expired and drops.
func writeNetscape(w io.Writer, cookies []*kooky.Cookie) {
	fmt.Fprint(w, "# Netscape HTTP Cookie File\n\n")
	for _, cookie := range cookies {
		domain := cookie.Domain
		if cookie.HttpOnly {
			domain = "#HttpOnly_" + domain
		}
		path := cookie.Path
		if path == "" {
			path = "/"
		}
		var expires int64 // 0 means a session cookie
		if !cookie.Expires.IsZero() {
			expires = cookie.Expires.Unix()
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\t%s\n",
			domain, netscapeBool(strings.HasPrefix(cookie.Domain, ".")), path, netscapeBool(cookie.Secure),
			expires, cookie.Name, cookie.Value)
	}
}

func netscapeBool(b bool) string {
	if b {
		return "TRUE"
	}
	return "FALSE"
}

// writeJSON writes v (a cookie or a slice of cookies) as JSON, followed by a newline.
func writeJSON(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}
