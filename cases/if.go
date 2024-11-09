package main

func function() interface{} {
	return nil
}

func main() {
	const a = 10

	// Basic conditional
	if a > 5 {
		// code to execute if condition is true
	}

	// Conditional with else clause
	if a <= 7 && a > 0 {
		// code to execute if condition is true
	} else {
		// code to execute if condition is false
	}

	// Nested conditional
	if a > 2 {
		if a > 5 {
			// code to execute if both condition1 and condition2 are true
		}
	}

	// Nested conditional with else clause
	if a-1 > 10 {
		// code to execute if condition1 is true
	} else if a == 10 {
		// code to execute if condition2 is true
	} else {
		// code to execute if none of the above conditions are true
	}

	// Conditional with initializer
	if err := function(); err != nil {
		// handle error
	}
}
