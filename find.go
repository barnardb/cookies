package main

import (
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/zellyn/kooky"
	_ "github.com/zellyn/kooky/allbrowsers"
)

func storesForBrowsers(names []string) []kooky.CookieStore {
	stores := kooky.FindAllCookieStores()
	if names == nil {
		return stores
	}
	n := 0
STORES:
	for _, store := range stores {
		browser := store.Browser()
		for _, name := range names {
			if browser == name {
				stores[n] = store
				n++
				continue STORES
			}
		}
	}
	return stores[:n]
}

func findCookies(url *url.URL, browsers []string, logger *log.Logger) (cookies []*kooky.Cookie) {
	logger.Printf("Looking for cookies for URL %s", url)

	stores := storesForBrowsers(browsers)
	logger.Printf("Found %v cookie stores", len(stores))

	for _, store := range stores {
		logger.Printf("Loading cookies from %v", store)
		cookies, err := store.ReadCookies()
		if err != nil {
			logger.Printf("Error loading cookies from %v: %s", store, err)
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
	filter := currentlyAppliesToURL(url, logger)
	for _, cookie := range cookies {
		if filter(cookie) {
			filtered = append(filtered, cookie)
		}
	}
	logger.Printf("Accepted %d of %d cookies", len(filtered), len(cookies))
	return filtered
}

func currentlyAppliesToURL(url *url.URL, logger *log.Logger) kooky.Filter {
	currentTime := time.Now()
	logger.Printf("Current time is %v", currentTime)
	return appliesToURLAtTime(url, currentTime, logger)
}

func appliesToURLAtTime(url *url.URL, time time.Time, logger *log.Logger) kooky.Filter {
	isHttps := url.Scheme != "https"
	return func(cookie *kooky.Cookie) bool {
		if !hostMatchesDomain(url.Host, cookie.Domain) {
			logger.Printf("Rejecting cookie for non-matching domain: %v", cookie)
		} else if isHttps && cookie.Secure {
			logger.Printf("Rejecting secure cookie for non-HTTPS URL: %v", cookie)
		} else if !(cookie.Expires.IsZero() || time.Before(cookie.Expires)) {
			logger.Printf("Rejecting expired cookie: %v", cookie)
		} else if url.Path != "" && !strings.HasPrefix(url.Path, cookie.Path) {
			logger.Printf("Rejecting cookie due to unmatched path: %v", cookie)
		} else {
			logger.Printf("Accepting: %v", cookie)
			return true
		}
		return false
	}
}

func hostMatchesDomain(host string, domain string) bool {
	return host == domain ||
		(strings.HasPrefix(domain, ".") && (strings.HasSuffix(host, domain) || host == domain[1:]))
}
