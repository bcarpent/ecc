package main

import (
    "fmt"
    "math/big"
	"math/rand"
	"time"
)

// Pollard's Rho iterator type
type PollardRhoIterator struct {
	point1 *Point
	point2 *Point
	X1 *Point
	X2 *Point
	a1 int64
	b1 int64
	a2 int64
	b2 int64
	X *Point
	a int64
	b int64
}

func NewPollardRhoIterator(P *Point, Q *Point) (*PollardRhoIterator) {

	// Create new iterator
	iter := new(PollardRhoIterator)
	iter.point1 = P
	iter.point2 = Q

	// Generate random pair a1, b1 within range of n
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    iter.a1 = r1.Int63n(myCurve.n)
    if (iter.a1 == 0) {
    	iter.a1++
    }
    s2 := rand.NewSource(time.Now().UnixNano())
    r2 := rand.New(s2)
    iter.b1 = r2.Int63n(myCurve.n)
    if (iter.b1 == 0) {
    	iter.b1++
    }

    // TBD
//    iter.a1 = 1
//    iter.b1 = 3
//    fmt.Printf("a1 = %d, b1 = %d\n", iter.a1, iter.b1)

    // Compute a1P + b1Q
    a1P := Multiply(iter.a1, P)
//    fmt.Printf("a1P = (%d, %d)\n\n\n", a1P.x, a1P.y)
    b1Q := Multiply(iter.b1, Q)
 //   fmt.Printf("b1Q = (%d, %d)\n\n\n", b1Q.x, b1Q.y)
    iter.X1 = Add(a1P, b1Q)
/*    if (iter.X1 == nil) {
    	fmt.Printf("X1 = nil\n\n\n")
    } else {
	    fmt.Printf("X1 = (%d, %d)\n\n\n", iter.X1.x, iter.X1.y)
	}
*/
    // Generate random pair a2, b2 within range of n
    s1 = rand.NewSource(time.Now().UnixNano())
    r1 = rand.New(s1)
    iter.a2 = r1.Int63n(myCurve.n)
    if (iter.a2 == 0) {
    	iter.a2++
    }
    s2 = rand.NewSource(time.Now().UnixNano())
    r2 = rand.New(s2)
    iter.b2 = r2.Int63n(myCurve.n)
    if (iter.b2 == 0) {
    	iter.b2++
    }

     // Compute a2P + b2Q
    a2P := Multiply(iter.a2, P)
//    fmt.Printf("a2P = (%d, %d)\n", a2P.x, a2P.y)
    b2Q := Multiply(iter.b2, Q)
//    fmt.Printf("b2Q = (%d, %d)\n", b2Q.x, b2Q.y)


    iter.X2 = Add(a2P, b2Q)
    // Initialize iterated values X, a, b 
    iter.X = nil
    iter.a = 0
    iter.b = 0

	return iter
}

func (iter *PollardRhoIterator) Next() (X *Point, a int64, b int64) {
	var i *big.Int

    // Partition the curve into three segments
    partitionSize := (myCurve.p / 3) + 1
    bigPartitionSize := big.NewInt(partitionSize)

    if (iter.X == nil) {
		i = big.NewInt(0)
    } else {
    	i = new(big.Int)
    	i.Div(iter.X.x, bigPartitionSize)
    }

    // Start with 0P, then add 1P to get 1P, where i == 0.
    // Next iteration is 2P (doubling 1P to get 2P), where i == 1.
    // Next iteration is 3P (adding P to 2P), where i == 2.
    if (i.Int64() == 0) {
    	iter.a += iter.a1
    	iter.b += iter.b1
//    	fmt.Printf("Iterating, i = 0, a = %d, b = %d\n", iter.a, iter.b)
    	iter.X = Add(iter.X, iter.X1)
    } else if (i.Int64() == 1) {
    	iter.a *= 2
    	iter.b *= 2
//    	fmt.Printf("Iterating, i = 1, a = %d, b = %d\n", iter.a, iter.b)
    	iter.X = Double(iter.X)
    } else if (i.Int64() == 2) {
    	iter.a += iter.a2
    	iter.b += iter.b2
//    	fmt.Printf("Iterating, i = 2, a = %d, b = %d\n", iter.a, iter.b)
    	iter.X = Add(iter.X, iter.X2)
    } else {
//    	fmt.Printf("Invalid i, returning.....\n")
    	return nil, 0, 0
    }

    // Take the a, b values mod n
    iter.a = mod(iter.a, myCurve.n)
    iter.b = mod(iter.b, myCurve.n)

    a = iter.a
    b = iter.b

//    fmt.Printf("a, b = %d, %d\n\n\n", a, b)

    return iter.X, a, b
}

// Pollard's Rho algorithm
func pollardRho(P *Point, Q *Point) (log int64, numSteps int64) {
	var i int64

	// Ensure points P and Q are both on the curve
	if ((isOnCurve(*P) == false) || (isOnCurve(*Q) == false)) {
		fmt.Printf("Exiting, point P or Q is not on curve!\n")
		return 0, numSteps
	}

	// Initialize numSteps to 0
	numSteps = 0

	tortoise := NewPollardRhoIterator(P, Q)
	if (tortoise == nil) {
		return 0, 0
	}

	// Copy values of tortoise to the hare iterator
	hare := *tortoise

	for i = 0; i < myCurve.n; i++ {
		X1, a1, b1 := tortoise.Next()

		// Hare skips over a step every time
		X2, a2, b2 := hare.Next()
		X2, a2, b2 = hare.Next()

		numSteps++

		// If (x1,y1) == (x2,y2), we've found a match (detected a cycle)
		if ((X1 == nil) && (X2 == nil)) {
			if (b1 == b2) {
				fmt.Printf("Generated random sequence divide by zero, retry\n")
				return 0, numSteps
			}
			log = (a1 - a2) * modInverse((b2 - b1), myCurve.n)
			log = mod(log, myCurve.n)
			return log, numSteps
		} else if ((X1 == nil) || (X2 == nil)) {
			numSteps++
		} else if ((X1.x.Cmp(X2.x) == 0) && (X1.y.Cmp(X2.y) == 0)) {
			if (b1 == b2) {
				fmt.Printf("Generated random sequence divide by zero, retry\n")
				return 0, numSteps
			}
			log = (a1 - a2) * modInverse((b2 - b1), myCurve.n)
			log = mod(log, myCurve.n)
			return log, numSteps
		} else {
			continue
		}
	}

	return 0, numSteps
}
