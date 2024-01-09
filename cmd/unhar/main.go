package main

import (
	"fmt"
	"os"

	unhar "github.com/code-lion-com/go-unhar"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr,
			"Err: Missing filename(s)\n",
			"usage: unhar [filename] [filename] ...",
		)
		os.Exit(1)
	}

	for _, filePath := range os.Args[1:] {
		har := &unhar.Har{}
		err := har.Open(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reading %q: %s\n", filePath, err)
			os.Exit(1)
		}

		err = har.Write(filePath, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "extracting %q: %s\n", filePath, err)
			os.Exit(1)
		}
	}
}
