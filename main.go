package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	flag "github.com/spf13/pflag"
)

type options struct {
	browsers      []string
	url           *url.URL
	acceptMissing bool
	verbose       bool
}

func main() {
	options := parseCommandLine()
	logger := buildLogger(options.verbose)

	cookies := findCookies(options.url, options.browsers, logger)
	if len(cookies) == 0 {
		if !options.acceptMissing {
			os.Exit(1)
		}
		return
	}

	formatCookies(os.Stdout, cookies)
	fmt.Print("\n")
}

func parseCommandLine() (options options) {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	fatalError := func(error ...interface{}) {
		fmt.Fprintln(os.Stderr, error...)
		flagSet.Usage()
		os.Exit(2)
	}

	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options…] <URL>\n\nThe following options are available:\n", os.Args[0])
		flagSet.PrintDefaults()
	}

	flagSet.BoolVarP(&options.acceptMissing, "accept-missing", "a", false, "don't fail with exit status 1 when cookies aren't found")
	flagSet.StringArrayVarP(&options.browsers, "browser", "b", []string{"chrome", "chromium", "firefox", "safari"}, "browser to try extracting a cookie from, can be repeated to try multiple browsers")
	flagSet.BoolVarP(&options.verbose, "verbose", "v", false, "enables logging to stderr")

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		fatalError(err)
	}

	if flagSet.NArg() != 1 {
		fatalError("error: expected 1 argument but got", flag.NArg())
	}

	options.url, err = url.Parse(flagSet.Arg(0))
	if err != nil {
		fatalError("error parsing URL:", err)
	} else if options.url.Scheme != "http" && options.url.Scheme != "https" {
		fatalError(fmt.Sprintf("URL scheme must be http or https, but got \"%s\"", options.url.Scheme))
	} else if options.url.Host == "" {
		fatalError("URL host must be non-empty")
	}

	return
}

func buildLogger(verbose bool) *log.Logger {
	w := ioutil.Discard
	if verbose {
		w = os.Stderr
	}
	return log.New(w, "", 0)
}
