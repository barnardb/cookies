package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/all"
)

func contains(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

// profileMatcher selects cookie stores by profile name. A store is selected
// if its profile equals one of names or matches one of patterns. If both are
// empty, every profile is selected.
type profileMatcher struct {
	names    []string
	patterns []*regexp.Regexp
}

func (m profileMatcher) matches(profile string) bool {
	if len(m.names) == 0 && len(m.patterns) == 0 {
		return true
	}
	if contains(m.names, profile) {
		return true
	}
	for _, pattern := range m.patterns {
		if pattern.MatchString(profile) {
			return true
		}
	}
	return false
}

func findCookies(url *url.URL, name string, browsers []string, profiles profileMatcher, logger *Logger) (results []*kooky.Cookie) {
	logger.Printf("Looking for browsers %v", browsers)
	logger.Printf("Looking for cookies for URL %s", url)
	filter := currentlyAppliesToURLAndName(url, name, logger.RequireVerbosity(2))

	cookies := make(chan *kooky.Cookie)
	go func() {
		var wg sync.WaitGroup
		for store, err := range kooky.TraverseCookieStores(context.Background()) {
			if err != nil {
				// Finders report every browser that isn't installed, so this is only interesting in detail.
				logger.RequireVerbosity(2).Printf("Error finding cookie stores: %s", err)
				continue
			}
			if !contains(browsers, store.Browser()) {
				store.Close()
				continue
			}
			if _, err := os.Stat(store.FilePath()); err != nil {
				// Finders also propose candidate paths that may not exist, e.g. Chrome's pre-v96 location.
				logger.RequireVerbosity(2).Printf("Skipping cookies from %s: %s", describeStore(store), err)
				store.Close()
				continue
			}
			if !profiles.matches(store.Profile()) {
				logger.Printf("Skipping cookies from %s: profile %q not selected", describeStore(store), store.Profile())
				store.Close()
				continue
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer store.Close()
				logger.Printf("Loading cookies from %s", describeStore(store))
				for cookie, err := range store.TraverseCookies(filter) {
					if err != nil {
						logger.Printf("Error loading cookies from %s: %s", describeStore(store), err)
						continue
					}
					cookies <- cookie
				}
				logger.Printf("Done loading cookies from %s", describeStore(store))
			}()
		}
		wg.Wait()
		close(cookies)
	}()
	for cookie := range cookies {
		results = append(results, cookie)
	}
	logger.Printf("Found %d matching cookie(s)", len(results))
	return
}

func describeStore(store kooky.CookieStore) string {
	return fmt.Sprintf("%s profile %q (%s)", store.Browser(), store.Profile(), store.FilePath())
}

func currentlyAppliesToURLAndName(url *url.URL, name string, logger *Logger) kooky.Filter {
	currentTime := time.Now()
	logger.Printf("Current time is %v", currentTime)
	return appliesToURLAndNameAtTime(url, name, currentTime, logger)
}

func appliesToURLAndNameAtTime(url *url.URL, name string, time time.Time, logger *Logger) kooky.Filter {
	urlIsNotSecure := url.Scheme != "https"
	// A plain FilterFunc (not a ValueFilterFunc) lets kooky apply the filter before decrypting
	// values, so we only touch the keychain for browsers that have a matching cookie.
	return kooky.FilterFunc(func(cookie *kooky.Cookie) bool {
		if !hostMatchesDomain(url.Host, cookie.Domain) {
			logger.Printf("Rejecting cookie for non-matching domain: %v", cookie)
		} else if urlIsNotSecure && cookie.Secure {
			logger.Printf("Rejecting secure cookie for non-HTTPS URL: %v", cookie)
		} else if !(cookie.Expires.IsZero() || time.Before(cookie.Expires)) {
			logger.Printf("Rejecting expired cookie: %v", cookie)
		} else if url.Path != "" && !strings.HasPrefix(url.Path, cookie.Path) {
			logger.Printf("Rejecting cookie due to unmatched path: %v", cookie)
		} else if name != "" && cookie.Name != name {
			logger.Printf("Rejecting cookie due to unmatched name: %v", cookie)
		} else {
			logger.Printf("Accepting: %v", cookie)
			return true
		}
		return false
	})
}

func hostMatchesDomain(host string, domain string) bool {
	return host == domain ||
		(strings.HasPrefix(domain, ".") && (strings.HasSuffix(host, domain) || host == domain[1:]))
}
