package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	input := bufio.NewReader(file)
	lexer := NewLexer(input)
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
