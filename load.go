package main

import (
  "fmt"
  "github.com/barnardb/kooky"
  "os/user"
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
    cookiesFile, err := pathFromHome("/Library/Application Support/Firefox/Profiles/cp9hpajj.default/cookies.sqlite")
    if err != nil {
      return nil, err
    }
    return kooky.ReadFirefoxCookies(cookiesFile)
}

func loadSafariCookies(domain string) ([]*kooky.Cookie, error) {
    cookiesFile, err := pathFromHome("/Library/Cookies/Cookies.binarycookies")
    if err != nil {
      return nil, err
    }
    return kooky.ReadSafariCookies(cookiesFile, domain, "", time.Time{})
}

func getCookieLoader(name string) (loader cookieLoader, err error) {
  normalized := strings.ToLower(name); switch normalized {
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
