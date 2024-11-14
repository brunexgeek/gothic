package main

import "fmt"

// Basic declaration with initialization
var x int = 10

// Declaration without initialization (zero value)
var y int

// Multiple variables of the same type
var a, b, c int

// Multiple variables with initialization
var d, e, f int = 1, 2, 3

// Mixed types with initialization
var (
	g1, g2 int    = 1, 1
	h      string = "hello"
	i      bool   = true
)

// Constants (const keyword)
const pi = 3
const (
	k = 1
	l = "constant string"
)

func main() {
	// Basic declaration with initialization
	var x int = 10
	fmt.Println(x)

	// Declaration without initialization (zero value)
	var y int
	fmt.Println(y)

	// Multiple variables of the same type
	var a, b, c int
	fmt.Println(a, b, c)

	// Multiple variables with initialization
	var d, e, f int = 1, 2, 3
	fmt.Println(d, e, f)

	// Mixed types with initialization
	var (
		g        = 1
		h string = "hello"
		i bool   = true
	)
	fmt.Println(g, h, i)

	// Constants (const keyword)
	const pi = 3.1415
	const (
		k int = 1
		l     = "constant string"
	)

	// Short variable declaration
	j := 10
	fmt.Println(j)

	// Declaration with the new function
	m := new(int)
	fmt.Println(*m) // Outputs: 0

	// Declaration with the make function
	n := make([]int, 10)
	o := make(map[string]int)
	p := make(chan int)
	fmt.Println(n, o, p)
}
