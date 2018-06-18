package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("err", err)
		os.Exit(1)
	}
	fmt.Println(string(body))
}
