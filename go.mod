module github.com/barnardb/cookies

go 1.23.10

toolchain go1.24.11

require (
	github.com/spf13/pflag v1.0.6
	github.com/zellyn/kooky v0.0.0-20201108220156-bec09c12c339
)

require (
	al.essio.dev/pkg/shellescape v1.5.1 // indirect
	github.com/Velocidex/json v0.0.0-20220224052537-92f3c0326e5a // indirect
	github.com/Velocidex/ordereddict v0.0.0-20250626035939-2f7f022fc719 // indirect
	github.com/Velocidex/yaml/v2 v2.2.8 // indirect
	github.com/alecthomas/assert v1.0.0 // indirect
	github.com/alecthomas/colour v0.1.0 // indirect
	github.com/alecthomas/repr v0.1.1 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/bobesa/go-domain-util v0.0.0-20190911083921-4033b5f7dd89 // indirect
	github.com/browserutils/kooky v0.2.4 // indirect
	github.com/creack/pty v1.1.9 // indirect
	github.com/danieljoos/wincred v1.2.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-sqlite/sqlite3 v0.0.0-20180313105335-53dd8e640ee7 // indirect
	github.com/godbus/dbus v4.1.0+incompatible // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gonuts/binary v0.2.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/keybase/dbus v0.0.0-20220506165403-5aa21ea2c23a // indirect
	github.com/keybase/go-keychain v0.0.1 // indirect
	github.com/keybase/go.dbus v0.0.0-20200324223359-a94be52c0b03 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/kr/pty v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sebdah/goldie v1.0.0 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/smartystreets/assertions v0.0.0-20180927180507-b2de0cb4f26d // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/yuin/goldmark v1.4.13 // indirect
	github.com/zalando/go-keyring v0.2.6 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/mod v0.25.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/telemetry v0.0.0-20240521205824-bda55230c457 // indirect
	golang.org/x/term v0.33.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	golang.org/x/tools v0.34.0 // indirect
	golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	www.velocidex.com/golang/go-ese v0.2.0 // indirect
)

// We are currently using the head of <https://github.com/barnardb/kooky/tree/visitors>.
// I decided to use my fork while waiting for https://github.com/zellyn/kooky/pull/43
// to be merged, but then decided to do some work to allow the kind of
// parallelism and defered decrypting I wanted to have. The result is a
// speed-up taking a typical use case of mine from ~1200ms to ~100ms.
// I intend to try to contribute these changes back upstream.
replace github.com/zellyn/kooky => github.com/barnardb/kooky v0.0.0-20210125015714-b0182bf77e67
