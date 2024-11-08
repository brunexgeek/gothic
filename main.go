package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	input, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	lexer := NewLexer(string(input))
	for {
		token, err := lexer.NextToken()
		fmt.Println(token)
		if err != nil {
			fmt.Println(err)
			break
		}
		if token.Type == EOF {
			break
		}
	}
}
