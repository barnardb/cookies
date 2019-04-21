package main

import (
	"github.com/barnardb/kooky"
	"log"
	"net/url"
	"strings"
	"time"
)

func findCookies(url *url.URL, browsers []string, logger *log.Logger) (cookies []*kooky.Cookie) {
	logger.Printf("Looking for cookies for URL %s", url)

	for _, browser := range browsers {
		loader, err := getCookieLoader(browser)
		if err != nil {
			logger.Printf("Error getting cookie loader for %s: %s", browser, err)
			continue
		}

		logger.Printf("Loading cookies from %s", browser)
		cookies, err := loader(url.Host)
		if err != nil {
			logger.Printf("Error loading cookies from %s: %s", browser, err)
			continue
		}
		cookies = filterCookies(cookies, url, logger)

		if len(cookies) > 0 {
			return cookies
		}
	}

	return []*kooky.Cookie{}
}

func filterCookies(cookies []*kooky.Cookie, url *url.URL, logger *log.Logger) []*kooky.Cookie {
	logger.Printf("Filtering %d cookies", len(cookies))
	filtered := []*kooky.Cookie{}
	now := time.Now()
	logger.Printf("Current time is %v", now)
	for _, cookie := range cookies {
		if cookie.Domain != url.Host {
			logger.Printf("Rejecting cookie for non-matching domain: %v", cookie)
		} else if url.Scheme != "https" && cookie.Secure {
			logger.Printf("Rejecting secure cookie for non-HTTPS URL: %v", cookie)
		} else if !cookie.Expires.IsZero() && cookie.Expires.Before(now) {
			logger.Printf("Rejecting expired cookie: %v", cookie)
		} else if url.Path != "" && !strings.HasPrefix(url.Path, cookie.Path) {
			logger.Printf("Rejecting cookie due to unmatched path: %v", cookie)
		} else {
			logger.Printf("Accepting: %v", cookie)
			filtered = append(filtered, cookie)
		}
	}
	logger.Printf("Accepted %d of %d cookies", len(filtered), len(cookies))
	return filtered
}
