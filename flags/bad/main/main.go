package main

import (
	"flag"
	"fmt"

	_ "playground/flags/bad"
)

func main() {
	v := flag.String("some", "", "Some flag, defined in main.go")
	flag.Parse()
	fmt.Println(*v)
}
