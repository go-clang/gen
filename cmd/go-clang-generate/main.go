package main

import (
	"fmt"
	"os"
)

func main() {
	err := Cmd(os.Args[1:])
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}
