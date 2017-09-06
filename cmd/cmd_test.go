package cmd_test

import (
	"fmt"

	"github.com/lestrrat/go-scripting/cmd"
)

func ExampleCommand() {
	res, err := cmd.New("ls", "-l").
		CaptureStdout(true).
		Grep(`_test\.go$`).
		Do(nil)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	fmt.Printf("%s\n", res.OutputString())
	// OUTPUT:
}
