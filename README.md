Cookies
=======

Extracts cookies from the users Chrome, Firefox or Safari cookie database,
outputting them in a format appropriate for use in the HTTP `Cookie` header.
This is useful in some scripting situations.

This core cookie reading code is provided by the [zellyn/kooky] cookie extraction library
This project wraps that library with some code to abstract browser differences away,
filter for cookies that match a URL, and provide a command-line interface.

[zellyn/kooky]: https://github.com/zellyn/kooky


- [Installing](#installing)
- [Usage](#usage)
    - [cURL Example](#curl-example)
    - [HTTPie Example](#httpie-example)
- [Status](#status)
- [Building](#building)


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

I use this tool on multiple MacOS systems.
(Note that there is a permissions issue you will have to deal with if you want to read Safari cookies: [zellyn/kooky #7].)
It should also work on other systems.
Pull requests are welcome.

[zellyn/kooky #7]: https://github.com/zellyn/kooky/issues/7


Building
--------

Requires [Go](https://golang.org/). Known to work with version `go1.15.6`.

Check out the repository and run
```bash
go build
```

This produces a `cookies` executable.
