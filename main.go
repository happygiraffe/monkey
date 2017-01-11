package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello, %s!  This is Monkey.\n", u.Username)

	repl.Start(os.Stdin, os.Stdout)
}
