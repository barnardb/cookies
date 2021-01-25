module github.com/barnardb/cookies

go 1.15

require (
	github.com/spf13/pflag v1.0.5
	github.com/zellyn/kooky v0.0.0-20201108220156-bec09c12c339
)

// We are currently using the head of <https://github.com/barnardb/kooky/tree/visitors>.
// I decided to use my fork while waiting for https://github.com/zellyn/kooky/pull/43
// to be merged, but then decided to do some work to allow the kind of
// parallelism and defered decrypting I wanted to have. The result is a
// speed-up taking a typical use case of mine from ~1200ms to ~100ms.
// I intend to try to contribute these changes back upstream.
replace github.com/zellyn/kooky => github.com/barnardb/kooky v0.0.0-20210125015714-b0182bf77e67
