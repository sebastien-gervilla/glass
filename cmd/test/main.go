package main

import (
	"fmt"
	"glass/language/repl"
	"os"
)

func main() {
	fmt.Print("Testing the glass programming language : \n")
	repl.Start(os.Stdin, os.Stdout)
}
