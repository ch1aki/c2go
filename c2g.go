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
	var max uint64
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

	lines := make([]Line, 0, 1024)
	if err := Parser(os.Stdin, &lines); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, i := range lines {
		if max < i.Count {
			max = i.Count
		}
	}

	for _, line := range lines {
		bar_count := int(line.Count * 30 / max)
		fmt.Printf("%s%d [%-30s] %s\n",
			line.Spacer, line.Count, strings.Repeat("|", bar_count), line.Text)
	}

	return ExitCodeOK
}

func Parser(stdin io.Reader, lines *[]Line) error {
	scanner := bufio.NewScanner(stdin)
	rep := regexp.MustCompile(`^(\s*)([0-9]+) (.*)$`)
	for scanner.Scan() {
		line := rep.FindSubmatch([]byte(scanner.Text()))
		spacer := string(line[1])
		count, err := strconv.ParseUint(string(line[2]), 10, 64)
		if err != nil {
			return fmt.Errorf("%s", "Parse Error")
		}
		text := string(line[3])

		*lines = append(*lines, Line{Spacer: spacer, Count: count, Text: text})
	}

	return nil
}
