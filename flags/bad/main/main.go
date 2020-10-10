package main

import (
	"flag"
	"fmt"

	_ "playground/flags/bad"
)

func main() {
	v := flag.String("some-flag", "", "Some flag, defined in main.go")
	flag.Parse()
	fmt.Println(*v)
}
