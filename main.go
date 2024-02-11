package main

import (
	"fmt"

	"github.com/dworthen/goscripty/cmd"
)

func main() {
	test("cool")
	cmd.Execute()
}

func test(name string) {
	fmt.Println(name)
}
