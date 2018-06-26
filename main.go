package main

import (
	"os"
)

const Version string = "v0.1.0"

func main() {
	cli := &CLI{inStream: os.Stdin, outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}
