Cookies
=======

Extracts cookies from the users Chrome, Firefox or Safari cookie database,
outputting them in a format appropriate for use in the HTTP `Cookie` header.
This is useful in some scripting situations.

This core cookie reading code is provided by the [zellyn/kooky] cookie extraction library.
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


Usage
-----

As explained by `cookies --help`:
```text
usage: cookies [options…] <URL> [<cookie-name>]

The following options are available:
  -a, --accept-missing        don't fail with exit status 1 when cookies aren't found
  -b, --browser stringArray   browser to try extracting a cookie from, can be repeated to try multiple browsers (default [chrome,chromium,firefox,safari])
  -v, --verbose               enables logging to stderr
```

So you get all cookies for a URL, so e.g. this:
```bash
cookies http://www.example.com
``` 
might yield
```
some.random.value=1234;JSESSIONID=0123456789ABCDEF0123456789ABCDEF;another_cookie:example-cookie-value
```

Or you can get a particular cookie value, so e.g. this:
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
