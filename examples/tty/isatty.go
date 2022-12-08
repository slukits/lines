package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
)

func main() {
	fmt.Println(isatty.IsTerminal(os.Stdout.Fd()))
	term := os.Getenv("TERM")
	fmt.Println("term:", term)
}
