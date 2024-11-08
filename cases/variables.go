package main

import "fmt"

// Basic Declaration with Initialization
var x int = 10

// Declaration Without Initialization (Zero Value)
var y int

// Multiple Variables of the Same Type
var a, b, c int

// Multiple Variables with Initialization
var d, e, f int = 1, 2, 3

// Mixed Types with Initialization
var (
	g int    = 1
	h string = "hello"
	i bool   = true
)

// Constants (const Keyword)
const pi = 3.14
const (
	k = 1
	l = "constant string"
)

func main() {
	// Basic Declaration with Initialization
	var x int = 10
	fmt.Println(x)

	// Declaration Without Initialization (Zero Value)
	var y int
	fmt.Println(y)

	// Multiple Variables of the Same Type
	var a, b, c int
	fmt.Println(a, b, c)

	// Multiple Variables with Initialization
	var d, e, f int = 1, 2, 3
	fmt.Println(d, e, f)

	// Mixed Types with Initialization
	var (
		g        = 1
		h string = "hello"
		i bool   = true
	)
	fmt.Println(g, h, i)

	// Constants (const Keyword)
	const pi = 3.1415
	const (
		k int = 1
		l     = "constant string"
	)

	// Short Variable Declaration
	j := 10
	fmt.Println(j)

	// Declaration with the new Function
	m := new(int)
	fmt.Println(*m) // Outputs: 0

	// Declaration with the make Function
	n := make([]int, 10)
	o := make(map[string]int)
	p := make(chan int)
	fmt.Println(n, o, p)
}
