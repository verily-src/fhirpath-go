package stablerand

import (
	"math/rand"
	"sync"
	"time"
)

var (
	// stableRand is the random engine used for generating random data in package
	// fhirtest.
	stableRand *rand.Rand

	// randMutex provides thread-safety for stableRand, in case any tests are
	// executed with t.Parallel(). Parallelism will affect the stability of the
	// randomness, since the generated values will no longer be deterministic once
	// concurrency is involved; but this doesn't mean the code should fail.
	randMutex sync.Mutex
)

const (
	// randSeed is the seed used for the random engine used in package fhirtest.
	// This seed is constant so that subsequent test executions will always receive
	// the same data.
	randSeed = 0xbadc0ffee

	// alnumAlphabet is a string containing all the upper and lowercase ascii
	// characters for letters and digits.
	alnumAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// hexAlphabet is a string containing lowercase hex characters.
	hexAlphabet = "abcdef0123456789"

	// decAlphabet is a string containing all decimal ascii characters.
	decAlphabet = "0123456789"
)

func init() {
	// Seed the random engine with a static value so that generation is consistent
	// across test executions, but still produces "unique" values.
	stableRand = rand.New(rand.NewSource(randSeed))
}

// Intn returns, as an int, a non-negative pseudo-random number in the half-open
// interval [0,n). It panics if n <= 0.
func Intn(n int) int {
	randMutex.Lock()
	defer randMutex.Unlock()

	return stableRand.Intn(n)
}

// Int63n returns, as an int64, a non-negative pseudo-random number in the
// half-open interval [0,n). It panics if n <= 0.
func Int63n(n int64) int64 {
	randMutex.Lock()
	defer randMutex.Unlock()

	return stableRand.Int63n(n)
}

// String returns, as a string, a pseudo-random string containing n characters
// all consisting of values within the supplied alphabet string.
// It panics if the alphabet string is empty.
func String(n int, alphabet string) string {
	if alphabet == "" {
		panic("No alphabet specified")
	}
	randMutex.Lock()
	defer randMutex.Unlock()
	b := make([]rune, n)
	for i := range b {
		b[i] = rune(alphabet[stableRand.Intn(len(alphabet))])
	}
	return string(b)
}

// AlnumString returns, as a string, a pseudo-random string containing n
// alphanumeric characters.
func AlnumString(n int) string {
	return String(n, alnumAlphabet)
}

// HexString returns, as a string, a pseudo-random string containing n
// hex characters.
func HexString(n int) string {
	return String(n, hexAlphabet)
}

// DecString returns, as a string, a pseudo-random string containing n
// decimal characters.
func DecString(n int) string {
	return String(n, decAlphabet)
}

// Time returns, as a time.Time object, a pseudo-random time starting with the
// base time, and adding a random amount between the half-open interval
// [0, variation) to the time. It panics if variation is negative.
func Time(base time.Time, variation time.Duration) time.Time {
	randMutex.Lock()
	defer randMutex.Unlock()

	offset := time.Duration(stableRand.Int63n(int64(variation)))
	base.Add(offset)
	return base
}

// OneOf returns, as a T object, a pseudo-randomly selected value from args.
// It panics if args is empty.
func OneOf[T any](args ...T) T {
	if len(args) == 0 {
		panic("No arguments specified to OneOf")
	}
	randMutex.Lock()
	defer randMutex.Unlock()

	i := stableRand.Intn(len(args))

	return args[i]
}
