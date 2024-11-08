package main

import "fmt"

func main() {
	// C-style loop
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}

	// Conditional loop
	i := 0
	for i < 10 {
		fmt.Println(i)
		i++
	}

	// Infinite loop
	for {
		fmt.Println("Looping forever")
		break // Use break to exit the loop
	}

	// Array and slice iteration with range
	array := []int{1, 2, 3, 4, 5}
	for index, value := range array {
		fmt.Printf("Index: %d, Value: %d\n", index, value)
	}

	// Map iteration with range
	m := map[string]int{"one": 1, "two": 2}
	for key, value := range m {
		fmt.Printf("Key: %s, Value: %d\n", key, value)
	}

	// String itaration with range
	s := "hello"
	for index, runeValue := range s {
		fmt.Printf("Index: %d, Rune: %c\n", index, runeValue)
	}

}
