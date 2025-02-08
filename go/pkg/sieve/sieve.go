package sieve

import (
	"math"
)

/*
Further Improvements:
	- Add wheel optimization to the basic sieve of eratosthenes
	- Implement Sieve of Atkin
	- Add concurrency to the Segmented Sieve
	- Implement config (and/or) parameters to choose which internal sieve to use
*/

// Sieve - provides an API for retrieving the Nth prime number using 0-based indexing where the 0th prime number is 2
type Sieve interface {
	NthPrime(n int64) int64
}

// PrimeNumberSieve - a struct required to implement the NthPrime Sieve interface.
type PrimeNumberSieve struct{}

// NewPrimeNumberSieve - Creates a new PrimeNumberSieve
func NewPrimeNumberSieve() *PrimeNumberSieve {
	return &PrimeNumberSieve{}
}

// NthPrime - Will calculate up to the nth prime number starting at 2
// if n is negative, the program will return 0
func (s *PrimeNumberSieve) NthPrime(nthPrime int64) int64 {

	if nthPrime < 0 {
		return 0
	}

	// use segmented sieve by default
	sieveFunc := &segmentedSieve{}

	// Pick a good upper bound: https://en.wikipedia.org/wiki/Prime_number_theorem
	upperBounds := nthPrime * (int64)(math.Log(float64(nthPrime)))
	if nthPrime < 6 {
		upperBounds = 20 // handles n <= 5 better since log is small for these
	}

	// Sieves till the upperbound and tests if the nth prime number can be found in the result
	// If not, scale upperbound and start again
	for {
		res := sieveFunc.sieve(upperBounds)
		if nthPrime <= int64(len(res)) {
			return res[nthPrime]
		}
		upperBounds = upperBounds * 2
	}
}

// sieve - internal interface used to switch between sieve implementations
// These functions are expected to return a list of primes from 2 - n.
// NOTE: This is not the same as the nth prime number.
type sieve interface {
	sieve(n int64) []int64
}

// segmentedSieve - uses a segmented sieve to return a list of primes from 2 - n.
// It starts by using the basic sieve of Erastothenes to return a list of primes from 2 - sqrt of n.
// Following that it creates segments to loop through, marking off any additional composites in the process
// finally, it adds the remaining primes before moving onto the next segment.
type segmentedSieve struct {
	basicSieve sieve
}

// sieve - implementation of the segmented sieve
func (s *segmentedSieve) sieve(n int64) []int64 {
	if s.basicSieve == nil {
		s.basicSieve = &basicSieveOfEratosthenes{}
	}

	// get segment size, use sqrt n as its consistent with what the basic sieve will use
	segmentSize := int64(math.Sqrt(float64(n)))

	// initialize primes up to sqrt(n) using the already created basic sieve of eratosthenes
	primes := s.basicSieve.sieve(segmentSize)

	// ensure the results contain all primes up to the square root of the upperbound found in the basic sieve
	result := make([]int64, 0)
	for _, p := range primes {
		result = append(result, p)
	}

	// begin processing segments
	for low := segmentSize; low <= n; low += segmentSize {

		high := low + segmentSize

		// high cannot be above the upperbound
		if high > n {
			high = n
		}

		// create a bool slice with enough capacity for the segment and mark them all as true
		segment := make([]bool, high-low+1)
		for i := range segment {
			segment[i] = true
		}

		for _, p := range primes {
			start := (low + p - 1) / p * p // find the smallest multiple of p that is greater than or equal to low
			// this is more performant than looping through and using % to find the start

			// ensure start is within the segment
			if start < low {
				start += p
			}

			// mark multiples of the prime as false in the segment
			for i := start; i <= high; i += p {
				segment[i-low] = false
			}
		}

		// Collect primes from the segment
		for i := low; i <= high; i++ {
			if segment[i-low] {
				result = append(result, i)
			}
		}
	}

	return result
}

// basicSieveOfEratosthenes - used to switch between sieve algorithms in the PrimeNumberSieve
type basicSieveOfEratosthenes struct{}

// basicSieveOfEratosthenes - uses a basic sieve of Erastothenes to return a list of primes from 2 - n
func (b *basicSieveOfEratosthenes) sieve(n int64) []int64 {

	// create a list of bools from 0 to upperbounds (n)
	isPrime := make([]bool, n+1)
	for i := 0; int64(i) <= n; i++ {
		isPrime[i] = true
	}

	// 0 and 1 are not prime numbers by definition so mark them false
	isPrime[0] = false
	isPrime[1] = false

	// loop through all primes from 2 to the square root of n (simple optimization: no need to check above sqrt(n) as a previous prime would already marked these)
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
