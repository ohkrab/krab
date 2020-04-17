package main

import (
	"fmt"
	"github.com/ohkrab/krab/krab"
	"github.com/ohkrab/krab/krab/parser"
)

func main() {
	fmt.Print("Krab v", krab.Version, "\n")

	parsed, err := parser.ParseFromFile("test/fixtures/migrations/create_table.hcl")
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println(parsed.Ast.Node)
}
