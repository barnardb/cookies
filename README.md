Cookies
=======

Extracts cookies from Chrome's cookie database,
outputting them in a format appropriate for use in the HTTP `Cookie` header.
This is useful in some scripting situations.

This core cookie reading code is provided by the [zellyn/kooky] cookie extraction library
(currently using a [fork] that merges [zellyn/kooky #5] to support newer Firefox databases).
This project wraps that library with some code to abstract browser differences away,
filter for cookies that match a URL, and provide a command-line interface.

[zellyn/kooky]: https://github.com/zellyn/kooky
[barnardb/kooky]: https://github.com/barnardb/kooky
[zellyn/kookie #5]: https://github.com/zellyn/kooky/pull/5


Status
------

The code is working on my MacOS 10.14.4 Mojave system,
though with an issue related to reading Safari cookies (see [zellyn/kooky #7]).
Paths to cookie files are currently MacOS-specific, but could easily be made OS-dependent.
Pull requests are welcome.

[zellyn/kooky #7]: https://github.com/zellyn/kooky/issues/7


Usage
-----

As explained by `cookies --help`:
```text
usage: ./cookies [options…] <URL>

The following options are available:
  -a, --accept-missing        don't fail with exit status 1 when cookies aren't found
  -b, --browser stringArray   browser to try extracting a cookie from, can be repeated to try multiple browsers (default [chrome,firefox,safari])
  -v, --verbose               enables logging to stderr
```

So for example:
```bash
cookies http://www.example.com
``` 
might yield
```
some.random.value=1234;JSESSIONID=0123456789ABCDEF0123456789ABCDEF;another_cookie:example-cookie-value
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

