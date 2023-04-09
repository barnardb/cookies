Cookies
=======

Extracts cookies from the user's Chrome, Chromium, Firefox or Safari cookie database.
A single cookie value can be retrieved, or all cookies applicable to given URL
can be retrieved and output in a format appropriate for use in the HTTP `Cookie` header.
Both of these usages are useful for scripting purposes.

The core cookie reading code is provided by the [zellyn/kooky] cookie extraction library.
This `cookies` tool provides a command-line interface to that library that
allows you to select which browser cookie databases to use and filter for all
or a particular cookie that are relevant for a given URL.

[zellyn/kooky]: https://github.com/zellyn/kooky


- [Installing](#installing)
- [Usage](#usage)
    - [cURL Example](#curl-example)
    - [HTTPie Example](#httpie-example)
- [Status](#status)
- [Building](#building)
- [Releasing](#releasing)


Installing
----------

On MacOS with Homebrew:

```bash
brew install barnardb/cookies/cookies
```

Alternatively, or on other platforms, follow the instructions below for [building](#building) an executable.


Usage
-----

As explained by `cookies --help`:
```text
usage: ./cookies [options…] <URL> [<cookie-name>]

The following options are available:
  -a, --accept-missing        don't fail with exit status 1 when cookies aren't found
  -b, --browser stringArray   browser to try extracting a cookie from, can be repeated to try multiple browsers (default [chrome,chromium,firefox,safari])
  -v, --verbose[=level]       enables logging to stderr; specify it twice or provide level 2 to get per-cookie details (`-vv` or `--verbose=2`)
      --version               prints version information and exits

cookies version 0.5.1  (https://github.com/barnardb/cookies)
```

To get all cookies relevant to a URL in the format expected by the `Cookie` header,
provide the URL as an argument. E.g., running
```bash
cookies http://www.example.com
```
might yield
```
some.random.value=1234;JSESSIONID=0123456789ABCDEF0123456789ABCDEF;another_cookie:example-cookie-value
```

Or you can get just the value of a particular cookie by providing both a URL and a cookie name.
E.g. running
```bash
cookies http://www.example.com JSESSIONID
```
might yield
```
0123456789ABCDEF0123456789ABCDEF
```

### cURL example

```bash
curl --cookie "$(cookies http://www.example.com)" http://www.example.com
```

might produce an HTTP request like this:

```http
GET / HTTP/1.1
Host: www.example.com
User-Agent: curl/7.54.0
Accept: */*
Cookie: some.random.value=1234;JSESSIONID=0123456789ABCDEF0123456789ABCDEF;another_cookie:example-cookie-value
```

### HTTPie example

```bash
http http://www.example.com Cookie:"$(cookies http://www.example.com)"
```

might produce an HTTP request like this:

```http
GET / HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Cookie: some.random.value=1234;JSESSIONID=0123456789ABCDEF0123456789ABCDEF;another_cookie:example-cookie-value
Host: www.example.com
User-Agent: HTTPie/1.0.2
```


Status
------

I use this tool for day-to-day tasks on multiple MacOS systems.
(Note that there is a cookie database permission issue you will have to deal with if you want to read Safari cookies: [zellyn/kooky #7].)

[zellyn/kooky #7]: https://github.com/zellyn/kooky/issues/7

As the library is essentially a wrapper around [zellyn/kooky] and the library
supports other platforms as well, this tool should also work on other platforms.

Pull requests are welcome.


Building
--------

Requires [Go](https://golang.org/).
Known to work with version `go1.15.6`.

To build the code, check out the repository and run:
```bash
go build
```

This produces a `cookies` executable.


Releasing
---------

Releases are prepared by running:
```bash
./prepare-release.sh "${version}"
```
`${version}` should be a semantic version number in the "0.0.0" format.
This tags the release (e.g. as "v0.0.0") and creates a draft release in GitHub,
which can be given release notes and published.

Once the new release it published to GitHub, the homebrew formula in
[barnardb/homebrew-cookies] should be updated following the instructions in
that repo's README.

[barnardb/homebrew-cookies]: https://github.com/barnardb/homebrew-cookies
