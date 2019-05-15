package main

import (
	"fmt"
	"github.com/zellyn/kooky"
	"io"
)

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
