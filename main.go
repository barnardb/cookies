package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/browserutils/kooky"
	flag "github.com/spf13/pflag"
)

// Overridden by the linker at build time.
var version string = "development"

type options struct {
	browsers      []string
	url           *url.URL
	name          string
	acceptMissing bool
	verbosity     int
	profiles      profileMatcher
	format        string
}

func main() {
	options := parseCommandLine()
	logger := LoggerWithVerbosity(options.verbosity)

	cookies := findCookies(options.url, options.name, options.browsers, options.profiles, logger)
	if len(cookies) == 0 {
		if !options.acceptMissing {
			os.Exit(1)
		}
		return
	}

	if len(options.name) > 0 {
		cookies = []*kooky.Cookie{strongestCookie(cookies)}
	}

	switch options.format {
	case formatHeader:
		if len(options.name) == 0 {
			formatCookies(os.Stdout, cookies)
		} else {
			fmt.Print(cookies[0].Value)
		}
		fmt.Print("\n")
	case formatNetscape:
		writeNetscape(os.Stdout, cookies)
	case formatJSON:
		var v any = cookies
		if len(options.name) > 0 {
			v = cookies[0]
		}
		if err := writeJSON(os.Stdout, v); err != nil {
			fmt.Fprintln(os.Stderr, "error writing JSON:", err)
			os.Exit(1)
		}
	}
}

func parseCommandLine() (options options) {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	usage := func(output io.Writer) {
		fmt.Fprintf(output, "usage: %s [options…] <URL> [<cookie-name>]\n\nThe following options are available:\n", os.Args[0])
		fmt.Fprint(output, regexp.MustCompile(`--verbose level  `).ReplaceAllString(flagSet.FlagUsages(), `--verbose[=level]`))
		fmt.Fprintf(output, "\ncookies version %s  (https://github.com/barnardb/cookies)\n", version)
	}

	fatalError := func(error ...interface{}) {
		fmt.Fprintln(os.Stderr, error...)
		usage(os.Stderr)
		os.Exit(2)
	}

	flagSet.Usage = func() { usage(os.Stdout) }

	flagSet.BoolVarP(&options.acceptMissing, "accept-missing", "a", false, "don't fail with exit status 1 when cookies aren't found")
	flagSet.StringArrayVarP(&options.browsers, "browser", "b", []string{"chrome", "chromium", "firefox", "safari"}, "browser to try extracting a cookie from, can be repeated to try multiple browsers")
	flagSet.StringVarP(&options.format, "format", "f", formatHeader, "output `format`: header (for the HTTP Cookie header, or just the value when a <cookie-name> is given), netscape (cookies.txt, for curl -b, wget --load-cookies or yt-dlp --cookies), or json")
	flagSet.CountVarP(&options.verbosity, "verbose", "v", "enables logging to stderr; specify it twice or provide `level` 2 to get per-cookie details (`-vv` or `--verbose=2`)")
	flagSet.StringArrayVarP(&options.profiles.names, "profile", "p", nil, "browser profile to extract cookies from, by exact name; can be repeated (default all profiles)")
	profileRegexes := flagSet.StringArray("profile-regex", nil, "browser profile to extract cookies from, by regular expression `pattern` (unanchored); can be repeated")

	versionFlag := flagSet.Bool("version", false, "prints version information and exits")

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		fatalError(err)
	}

	for _, pattern := range *profileRegexes {
		compiled, err := regexp.Compile(pattern)
		if err != nil {
			fatalError("error parsing --profile-regex:", err)
		}
		options.profiles.patterns = append(options.profiles.patterns, compiled)
	}

	if !contains(formats, options.format) {
		fatalError(fmt.Sprintf("--format must be one of %s, but got %q", strings.Join(formats, ", "), options.format))
	}

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
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

// strongestCookie returns the most specific of cookies, which must be non-empty.
func strongestCookie(cookies []*kooky.Cookie) *kooky.Cookie {
	strongest := cookies[0]
	for _, cookie := range cookies[1:] {
		if isStronger(cookie, strongest) {
			strongest = cookie
		}
	}
	return strongest
}

// isStronger reports whether a is more specific than b: a longer domain wins, then a host-only
// cookie (no leading dot) beats a domain cookie for the same domain, then a longer path wins,
// then the newer cookie.
func isStronger(a, b *kooky.Cookie) bool {
	aDomain, aIsDomainCookie := strings.CutPrefix(a.Domain, ".")
	bDomain, bIsDomainCookie := strings.CutPrefix(b.Domain, ".")
	if len(aDomain) != len(bDomain) {
		return len(aDomain) > len(bDomain)
	}
	if aIsDomainCookie != bIsDomainCookie {
		return !aIsDomainCookie
	}
	if len(a.Path) != len(b.Path) {
		return len(a.Path) > len(b.Path)
	}
	return a.Creation.After(b.Creation)
}
