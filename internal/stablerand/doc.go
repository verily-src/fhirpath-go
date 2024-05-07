/*
Package stablerand is a small helper utility that encapsulates its random engine
and always uses the same seed value for its randomness.

This ensures reproducibility and stability across executions, giving a
pseudo-random distribution, but with deterministic predictability. This is
primarily intended for generating content for tests, which ensures that inputs
are still pseudo-random, but predictible and consistent across unchanged
executions.

Functions in this package are thread-safe, although use in threaded contexts
will remove any guarantees of determinism.

Note: This is primarily used internally for the fhirtest package to implement
"random" IDs and meta-IDs so that test resources retain the same general
values across executions.
*/
package stablerand
