package main

import (
	"fmt"
	"math/big"
)

// Global variable myCurve initialized in main
var myCurve Curve

// Global big ints
var aBig *big.Int
var bBig *big.Int
var pBig *big.Int
var nBig *big.Int

// Elliptic Curve: y^2 = x^3 + ax + b
type Curve struct {
	a int64   // coefficient of x in equation
	b int64   // y-intercept in equation
	p int64   // order of the underlying field
	g Point   // (x,y) of the base point
	n int64   // order of the subgroup
}

type Point struct {
	x *big.Int
	y *big.Int
	keystring string
}

// Is the passed-in point on our curve?
func isOnCurve(point Point) bool {
	x := point.x
	y := point.y

	// Set global big ints
	aBig := big.NewInt(myCurve.a)
	bBig := big.NewInt(myCurve.b)
	pBig := big.NewInt(myCurve.p)

	// Compute left side and right side
	leftSide  := new(big.Int)
	rightSide := new(big.Int)

	// Temp variable t = ax
	t := new(big.Int).Mul(aBig, x)

	// y^2 = x^3 + ax + b
	leftSide.Mul(y, y)
	rightSide.Mul(x, x)
	rightSide.Mul(rightSide, x)
	rightSide.Add(rightSide, bBig)
	rightSide.Add(rightSide, t)

	// If the point is on the curve, then
	// y^2 - (x^3 + ax + b) mod p = 0
	test := new(big.Int)

	mod := test.Mod(t.Sub(leftSide, rightSide), pBig)
	if (mod.Int64() == 0) {
		return true
	} else {
		return false
	}
}

// This function returns the modulus p of an integer a
// such that a mod p = result. (Golang's % operator is
// a remainder and not the same as the modulus operator)
func mod(a int64, p int64) (result int64) {

	// Initialize big ints with the integers a and p
	aBig := big.NewInt(a)
	pBig := big.NewInt(p)

	// Use big package for modulo implementation
	mod := new(big.Int).Mod(aBig, pBig)

	return mod.Int64()
}

// This function returns the inverse of a modulo p, such that
// a*b mod p = 1. For example, 12^-1 mod 97  = 89 since
// (12)(89) mod 97 = 1
func modInverse(a int64, p int64) (result int64) {

	// Initialize big ints with the integers a and p
	aBig := big.NewInt(a)
	pBig := big.NewInt(p)

	fmt.Printf("aBig = %d, pBig = %d\n", aBig, pBig)

	// Use big package for inverse modulo implementation
	inverse := new(big.Int)
	inverse.ModInverse(aBig, pBig)
	fmt.Printf("Inverse = %d\n", inverse)

	return inverse.Int64()
}

func Add(point1 *Point, point2 *Point) (sum *Point) {
	// If point1 is nil, sum is point2
	// If point2 is nil, sum is point1
	if (point1 == nil) {
		return point2
	}
	if (point2 == nil) {
		return point1
	}

	if (!isOnCurve(*point1)) {
		fmt.Printf("Point %v not on curve\n", point1)
		return nil
	}

	if (!isOnCurve(*point2)) {
		fmt.Printf("Point %v not on curve\n", point2)
		return nil
	}

	if ((point1 == nil) || (point2 == nil)) {
		fmt.Printf("Point passed into add is nil\n")
		return nil
	}

	// Set global big ints
	aBig := big.NewInt(myCurve.a)
	pBig := big.NewInt(myCurve.p)

	x1 := point1.x
	y1 := point1.y
	x2 := point2.x
	y2 := point2.y

//	fmt.Printf("Adding (%d, %d) + (%d, %d)\n", x1, y1, x2, y2)

	if ((x1.Cmp(x2) == 0) && (y1.Cmp(y2) != 0)) {
		fmt.Printf("x1 == x2 but y1 != y2, return nil\n")
		return nil
	}

	// Now solve for slope. See p. 4 of accompanying document for equations.
	// If x1 equals x2, s is the slope of the tangent and we are point
	// doubling; otherwise, we are doing point addition.
	s := new(big.Int)
	modInv := new(big.Int)
	if (x1.Cmp(x2) == 0) {
		// s = (3 * x1 * x1 + myCurve.a) * modInverse(2 * y1, myCurve.p)
		three := new(big.Int).SetInt64(3)
		two   := new(big.Int).SetInt64(2)
		s.Mul(x1, x1)
		s.Mul(s, three)
		s.Add(s, aBig)
		modInv.Mul(two, y1)
		modInv.ModInverse(modInv, pBig)
		s.Mul(s, modInv)
	} else {
		// s = (y1 - y2) * modInverse(x1 - x2, myCurve.p)
		s.Sub(y1, y2)
		modInv.Sub(x1, x2)
		modInv.ModInverse(modInv, pBig)
		s.Mul(s, modInv)
	}

//	fmt.Printf("Slope = %d\n", s)

	// Now solve for the sum (x3, y3). See page 4 of accompanying document
	// for equations.
	// x3 := (s * s) - x1 - x2 
	x3 := new(big.Int)
	x3.Mul(s, s)
	x3.Sub(x3, x1)
	x3.Sub(x3, x2)

	// y3 := s * (x1 - x3) - y1
	y3 := new(big.Int)
	y3.Sub(x1, x3)
	y3.Mul(s, y3)
	y3.Sub(y3, y1)

//	fmt.Printf("Sum before mod = (%d, %d)\n", x3, y3)

	x3.Mod(x3, pBig)
	y3.Mod(y3, pBig)

	key := generateKey(x3, y3)
	result := Point{x3, y3, key}

	if !isOnCurve(result) {
		fmt.Printf("Sum (%d, %d) is not on curve\n", x3, y3)
		return nil
	}

//	fmt.Printf("Sum = (%d, %d)\n", x3, y3)

	return &result
}

func Double(point *Point) (result *Point) {
	return Add(point, point)
}

// Scalar multiply x * point using "double and add"
func Multiply(x int64, point *Point) (product *Point) {
	var result *Point

	if ((x == 0) || (point == nil)) {
		fmt.Printf("Null n or point passed into multiply\n")
		return nil
	}

	// Initialize sum to (0, 0)
	result = nil
	adder := point

	// POINT DOUBLE AND ADD
	// Continue double and add, shifting x right until 0.
	// If, for example, x = 9, then we start with x = 1001
	// in binary and proceed as follows:
	// 1001: 2P = adder, result = 1P
	// 100:  4P = adder, result = 1P
	// 10:   8P = adder, result = 1P
	// 1:    result = 8P + 1P = 9P
	for (x > 0) {
//		fmt.Printf("x = %d\n", x)
		if ((x & 1) == 1) {
			result = Add(result, adder)
			if (result == nil) {
				fmt.Printf("Exiting Multiply, point not on curve\n")
				return nil
			}
//			fmt.Printf("result = %v\n", result)
		}
		adder = Double(adder)
//		fmt.Printf("adder = %v\n", adder)
		x >>= 1
	}

//	fmt.Printf("Product = (%d, %d)\n", result.x, result.y)
	return result
}

