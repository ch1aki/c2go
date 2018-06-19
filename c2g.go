package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	rep := regexp.MustCompile(`^(\s*)([0-9]+) (.*)$`)
	for scanner.Scan() {
		line := rep.FindSubmatch([]byte(scanner.Text()))
		spacer := string(line[1])
		count, err := strconv.Atoi(string(line[2]))
		if err != nil {
			fmt.Println("err", err)
			os.Exit(1)
		}
		text := string(line[3])

		fmt.Println(spacer)
		fmt.Println(count)
		fmt.Println(text)
	}
}
