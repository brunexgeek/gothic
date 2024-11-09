package main

import (
	"encoding/json"
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
	/*for {
		token, err := lexer.NextToken()
		fmt.Println(token)
		if err != nil {
			fmt.Println(err)
			break
		}
		if token.Type == T_EOF {
			break
		}
	}*/
	parser := NewParser(lexer)
	result := parser.Parse()
	if len(parser.errors) > 0 {
		fmt.Println("Errors encountered during parsing:")
		for _, err := range parser.errors {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Parsing completed successfully")
		fmt.Print(json.MarshalIndent(result, "", "  "))
	}

}
