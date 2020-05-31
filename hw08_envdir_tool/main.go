package main

import (
	"fmt"
	"os"
)

const (
	minCountArgs = 3
)

func main() {
	if len(os.Args) < minCountArgs {
		fmt.Println("USE:\ngo-envdir <PATH_TO_ENV_DIR> <COMMAND> [arg...]")
		return
	}
	e, err := ReadDir(os.Args[1])
	if err != nil {
		panic(err)
	}
	os.Exit(RunCmd(os.Args[2:], e))
}
