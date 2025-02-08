package sieve

import (
	"math"
)

// Sieve - provides an API for retrieving the Nth prime number using 0-based indexing where the 0th prime number is 2
type Sieve interface {
	NthPrime(n int64) int64
}

// NewSieve - Creates a new Sieve
func NewSieve() *PrimeNumberSieve {
	return &PrimeNumberSieve{}
}

// PrimeNumberSieve - a struct required to implement the NthPrime Sieve interface.
type PrimeNumberSieve struct{}

// NthPrime - Will calculate up to the nth prime number starting at 2
// if n is negative, the program will return 0
func (s *PrimeNumberSieve) NthPrime(nthPrime int64) int64 {

	sieve := basicSieveOfEratosthenes{}

	if nthPrime < 0 {
		return 0
	}

	// Pick a good upper bound: https://en.wikipedia.org/wiki/Prime_number_theorem
	upperBounds := nthPrime * (int64)(math.Log(float64(nthPrime)))
	if nthPrime < 6 {
		upperBounds = 20 // handles n <= 5 better since log is small for these
	}

	// Sieves till the upperbound and tests if the nth prime number can be found in the result
	// If not, scale upperbound and start again
	for {
		res := sieve.sieve(upperBounds)
		if nthPrime <= int64(len(res)) {
			return res[nthPrime]
		}
		upperBounds = upperBounds * 2
	}
}

// basicSieveOfEratosthenes - will be needed to implement a new interface (not included in this commit) that can
// be used to switch between sieve algorithms in the PrimeNumberSieve
type basicSieveOfEratosthenes struct{}

// basicSieveOfEratosthenes - uses a basic sieve of Erastothenes to return a list of primes from 2 - N
func (b *basicSieveOfEratosthenes) sieve(n int64) []int64 {

	// create a list of bools from 0 to upperbounds (n)
	isPrime := make([]bool, n+1)
	for i := 0; int64(i) <= n; i++ {
		isPrime[i] = true
	}

	// 0 and 1 are not prime numbers by definition so mark them false
	isPrime[0] = false
	isPrime[1] = false

	// Loop through all primes from 2 to the square root of n (simple optimization: no need to check above sqrt(n) as a prime factor would already been found)
	// if i is marked as a prime, mark all multiples of i as composites (false)
	for i := 2; int64(i)*int64(i) <= n; i++ {
		if isPrime[i] {
			for j := i * i; int64(j) <= n; j += i {
				isPrime[j] = false
			}
		}
	}

	// append all primes from 2 to n to results and return
	res := make([]int64, 0)
	for i, potentialPrime := range isPrime {
		if potentialPrime {
			res = append(res, int64(i))
		}
	}

	return res
}
