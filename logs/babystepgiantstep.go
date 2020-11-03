package main

import (
    "fmt"
    "math/big"
	"math/rand"
	"time"
)

func generateKey(x *big.Int, y *big.Int) (key string) {
	return x.String() + ", " + y.String()
}

// Shanks' Baby-step giant step algorithm
func babyStepGiantStep(P *Point, Q *Point) (x int64, numSteps int64) {

	var K, jP *Point
	var i, j, m int64

	// Initialize numSteps to 0
	numSteps = 0

	// Initialize global big ints
	pBig := big.NewInt(myCurve.p)
	nBig := big.NewInt(myCurve.n)

	// Ensure points P and Q are both on the curve
	if ((isOnCurve(*P) == false) || (isOnCurve(*Q) == false)) {
		fmt.Printf("Exiting, point P or Q is not on curve!\n")
		return 0, numSteps
	}

	// Calculate m = sqrt(n) + 1
	// Initialize big int with the integer n
	mBig := new(big.Int).Sqrt(nBig)
	m = mBig.Int64()
	m++

	// The hash map is of size m with key-value pairs (jP, j). The table
	// is indexed by a keystring representation of point jP with the 
	// corresponding entry's value set to the integer j.
	hTable := make(map[string]int64)

	// Now pre-compute the baby steps (jP) and store them in a hash map,
	// starting with 0P + P = 1P.
	jP = nil
	for j = 1; j <= m; j++ {
		jP = Add(jP, P)
		if (jP != nil) {
			// Generate the key string as x and y concatenated strings
			keystring := generateKey(jP.x, jP.y)
			hTable[keystring] = j
		}
		numSteps++
	}

	// Now compute the giant steps (Q - imP) and check the hash table
	// for any matching point. Here we precompute S = -mP, and our
	// key K into the hash table starts with K = Q - 0mP
	K = Q
	negY := new(big.Int).Neg(P.y)
	modY := new(big.Int).Mod(negY, pBig)
	keystring := generateKey(P.x, modY)
	negP := Point{P.x, negY, keystring}

//	fmt.Printf("Multiply %d with (%d, %d)\n", m, negP.x, negP.y)
	S := Multiply(m, &negP)
	if (S == nil) {
		fmt.Printf("Exiting, point -mP is not on curve!\n")
		return 0, numSteps		
	}

	// Now lookup K = Q + iS, starting with K = Q + 0S,
	// where S = -mP
	for i = 0; i <= m; i++ {
		key := generateKey(K.x, K.y)
		j = hTable[key]		  // Starts with Q + 0S
		if (j > 0) {		  // Found a match
			x := j + (m * i)  // x = im + j
			numSteps = numSteps + i
			return x, numSteps
		} else {
			K = Add(K, S)	  // Did not find a match, K = K + S
			if (K == nil) {
				fmt.Printf("Exiting, point K is not on curve!\n")
				return 0, numSteps
			}
		}
	}

	fmt.Printf("Baby-step giant-step log not found\n")
	return 0, numSteps
}

func main() {
	var numSteps int64

	// Initialize base point and our curve
	xBig := new(big.Int).SetInt64(92)
	yBig := new(big.Int).SetInt64(207)
	key := generateKey(xBig, yBig)
	basePoint := Point{xBig, yBig, key}
	p := basePoint
	myCurve = Curve{1, 11, 709, basePoint, 727}
	fmt.Printf("My elliptic curve: y^2 = (x^3 + %dx + %d) mod %d\n",
				myCurve.a, myCurve.b, myCurve.p)
	fmt.Printf("Base point: (%d, %d), key %s\n",
				basePoint.x.Int64(), basePoint.y.Int64(), basePoint.keystring)
	fmt.Printf("Order of subgroup: %d\n", myCurve.n)

	// Base point must be on the curve
	if (isOnCurve(basePoint) == false) {
		fmt.Printf("Exiting, base point is not on curve!\n")
		return
	}

	// Use a random number x to compute point Q. Later we will
	// use this number x to check our result y from y = log(P, Q).
	// Give the random number generator a seed that changes
	// based on current time.
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    x := r1.Int63n(myCurve.n)
    if (x == 0) {
    	x++
    }

	fmt.Printf("Random x = %d\n", x)
	fmt.Printf("Point P: (%d, %d)\n", p.x, p.y)

	// Now compute Q from xP
	q := Multiply(x, &p)
	fmt.Printf("Point Q: (%d, %d)\n", q.x, q.y)

    fmt.Printf("\nBABY-STEP GIANT-STEP ALGORITHM:\n")
    fmt.Printf("-------------------------------\n")

	// Now use the baby-step, giant-step algorithm to determine
	// the discrete logarithm from points P and Q on the curve
	timeNow := time.Now().UnixNano()
	log, numSteps := babyStepGiantStep(&p, q)
	newTime := time.Now().UnixNano()
	timeBsgs := newTime - timeNow

	if (log != x) {
		fmt.Printf("Baby-step, giant-step did not correctly compute discrete logarithm\n")
		return
	}

	fmt.Printf("Baby-step giant step successful:\n")
	fmt.Printf("log(p, q) = %d\n", log)
	fmt.Printf("Number of steps = %d\n", numSteps)
	fmt.Printf("Execution time: %d ns\n", timeBsgs)

    fmt.Printf("\nPOLLARD'S RHO ALGORITHM:\n")
    fmt.Printf("--------------------------\n")
	timeNow = time.Now().UnixNano()
    log, numSteps = pollardRho(&p, q)
	newTime = time.Now().UnixNano()
	timePR := newTime - timeNow

    if (log != x) {
    	fmt.Printf("Pollard's Rho method did not correctly compute discrete logarithm\n")
    	return
    }

    fmt.Printf("Pollard's Rho method successful:\n")
	fmt.Printf("log(p, q) = %d\n", log)
	fmt.Printf("Number of steps = %d\n", numSteps + 1)
	fmt.Printf("Execution time: %d ns\n", timePR)

}

