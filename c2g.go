package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Line struct {
	Spacer string
	Count  uint64
	Text   string
}

func main() {
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

	fmt.Println(lines)
	fmt.Println(max)
}
