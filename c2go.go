package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

const (
	ExitCodeOK = iota
	ExitCodeParseError
	ExitCodeParseFlagError
)

type CLI struct {
	inStream             io.Reader
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
	var total uint64
	var totalFlag bool
	var target uint64
	var length int
	var verbose bool

	flags := flag.NewFlagSet("c2g", flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.BoolVar(&version, "version", false, "Print version information and quit")
	flags.BoolVar(&totalFlag, "total", false, "Print a bar graph of percentage in a relative to total.")
	flags.IntVar(&length, "length", 30, "Specify length of bar graph.")
	flags.BoolVar(&verbose, "verbose", false, "Verbose mode.")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagError
	}

	if version {
		fmt.Fprintf(c.errStream, "c2g version %s\n", Version)
		return ExitCodeOK
	}

	lines := make([]Line, 0, 1024)
	if err := Parser(c.inStream, &lines); err != nil {
		return ExitCodeParseError
	}

	for _, i := range lines {
		if max < i.Count {
			max = i.Count
		}
		total += i.Count
	}

	if totalFlag {
		target = total
	} else {
		target = max
	}

	switch {
	case length > 100:
		length = 100
	case length < 10:
		length = 10
	}

	PrintGraph(c.outStream, lines, target, length, verbose)

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

func PrintGraph(stdout io.Writer, lines []Line, target uint64, length int, verbose bool) {

	if verbose {
		for _, line := range lines {
			barCount := int(line.Count * uint64(length) / target)
			percentage := int(line.Count * 100 / target)
			var graph strings.Builder
			graph.WriteString(strings.Repeat("|", barCount))
			graph.WriteString(strings.Repeat(" ", length-barCount))

			fmt.Fprintf(stdout, "%s%d [%s%d] %s\n",
				line.Spacer, line.Count, graph.String()[:length-len(strconv.Itoa(percentage))], percentage, line.Text)
		}
	} else {
		format := fmt.Sprintf("%s%d%s", "%s%d [%-", length, "s] %s\n")
		for _, line := range lines {
			bar_count := int(line.Count * uint64(length) / target)
			fmt.Fprintf(stdout, format,
				line.Spacer, line.Count, strings.Repeat("|", bar_count), line.Text)
		}
	}
}
