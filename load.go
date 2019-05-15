package main

import (
	"fmt"
	"github.com/zellyn/kooky"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

type cookieLoader func(string) ([]*kooky.Cookie, error)

func pathFromHome(pathRootedAtHome string) (path string, err error) {
	usr, err := user.Current()
	if err != nil {
		return
	}
	return (usr.HomeDir + pathRootedAtHome), nil
}

func loadChromeCookies(domain string) ([]*kooky.Cookie, error) {
	cookiesFile, err := pathFromHome("/Library/Application Support/Google/Chrome/Default/Cookies")
	if err != nil {
		return nil, err
	}
	return kooky.ReadChromeCookies(cookiesFile, domain, "", time.Time{})
}

func loadFirefoxCookies() ([]*kooky.Cookie, error) {
	profilesDir, err := pathFromHome("/Library/Application Support/Firefox/Profiles")
	if err != nil {
		return nil, err
	}
	cookiesFiles, err := filepath.Glob(profilesDir + "/*.default/cookies.sqlite")
	if err != nil {
		return nil, err
	}
	if len(cookiesFiles) != 1 {
		return nil, fmt.Errorf("Expected to find one default Firefox profile with a cookies database, but found %d: %v", len(cookiesFiles), cookiesFiles)
	}
	return kooky.ReadFirefoxCookies(cookiesFiles[0])
}

func loadSafariCookies(domain string) ([]*kooky.Cookie, error) {
	cookiesFile, err := pathFromHome("/Library/Cookies/Cookies.binarycookies")
	if err != nil {
		return nil, err
	}
	return kooky.ReadSafariCookies(cookiesFile, domain, "", time.Time{})
}

func getCookieLoader(name string) (loader cookieLoader, err error) {
	normalized := strings.ToLower(name)
	switch normalized {
	case "chrome":
		loader = loadChromeCookies
	case "firefox":
		loader = func(_ string) ([]*kooky.Cookie, error) { return loadFirefoxCookies() }
	case "safari":
		loader = loadSafariCookies
	default:
		err = fmt.Errorf("No cookie loader matching %s", normalized)
	}
	return
}
