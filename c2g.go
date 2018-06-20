package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	ExitCodeOK = iota
	ExitCodeParseFlagError
)

type CLI struct {
	outStream, errStream io.Writer
}

type Line struct {
	Spacer string
	Count  uint64
	Text   string
}

func (c *CLI) Run(args []string) int {
	var version bool
	flags := flag.NewFlagSet("c2g", flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.BoolVar(&version, "version", false, "Print version information and quit")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagError
	}

	if version {
		fmt.Fprintf(c.errStream, "c2g version %s\n", Version)
		return ExitCodeOK
	}

	var max uint64
	lines := make([]Line, 0, 1024)
	scanner := bufio.NewScanner(os.Stdin)
	rep := regexp.MustCompile(`^(\s*)([0-9]+) (.*)$`)
	for scanner.Scan() {
		line := rep.FindSubmatch([]byte(scanner.Text()))
		spacer := string(line[1])
		count, err := strconv.ParseUint(string(line[2]), 10, 64)
		if err != nil {
			fmt.Println("err", err)
			os.Exit(1)
		}
		text := string(line[3])

		lines = append(lines, Line{Spacer: spacer, Count: count, Text: text})

		if count > max {
			max = count
		}
	}

	for _, line := range lines {
		bar_count := int(line.Count * 30 / max)
		fmt.Printf("%s%d [%-30s] %s\n",
			line.Spacer, line.Count, strings.Repeat("|", bar_count), line.Text)
	}

	return ExitCodeOK
}
