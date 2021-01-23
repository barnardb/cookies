package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/barnardb/kooky"
	flag "github.com/spf13/pflag"
)

type options struct {
	browsers      []string
	url           *url.URL
	name          string
	acceptMissing bool
	verbose       bool
}

func main() {
	options := parseCommandLine()
	logger := buildLogger(options.verbose)

	cookies := findCookies(options.url, options.name, options.browsers, logger)
	if len(cookies) == 0 {
		if !options.acceptMissing {
			os.Exit(1)
		}
		return
	}

	if len(options.name) == 0 {
		formatCookies(os.Stdout, cookies)
	} else {
		writeStrongestValue(os.Stdout, cookies)
	}
	fmt.Print("\n")
}

func parseCommandLine() (options options) {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	usage := func(output io.Writer) {
		fmt.Fprintf(output, "usage: %s [options…] <URL> [<cookie-name>]\n\nThe following options are available:\n", os.Args[0])
		fmt.Fprint(output, flagSet.FlagUsages())
	}

	fatalError := func(error ...interface{}) {
		fmt.Fprintln(os.Stderr, error...)
		usage(os.Stderr)
		os.Exit(2)
	}

	flagSet.Usage = func() { usage(os.Stdout) }

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

	if flagSet.NArg() != 1 && flagSet.NArg() != 2 {
		fatalError("error: expected 1 or 2 arguments but got", flag.NArg())
	}

	options.url, err = url.Parse(flagSet.Arg(0))
	if err != nil {
		fatalError("error parsing URL:", err)
	} else if options.url.Scheme != "http" && options.url.Scheme != "https" {
		fatalError(fmt.Sprintf("URL scheme must be http or https, but got \"%s\"", options.url.Scheme))
	} else if options.url.Host == "" {
		fatalError("URL host must be non-empty")
	}

	if flagSet.NArg() > 1 {
		options.name = flagSet.Arg(1)
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

func writeStrongestValue(w io.Writer, cookies []*kooky.Cookie) {
	strongest := cookies[0]
	for _, cookie := range cookies[1:] {
		if len(cookie.Domain) > len(strongest.Domain) || (cookie.Domain == strongest.Domain &&
			(len(cookie.Path) > len(strongest.Path) || (cookie.Path == strongest.Path &&
				cookie.Creation.After(strongest.Creation)))) {
			strongest = cookie
		}
	}
	fmt.Print(strongest.Value)
}
