package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type code struct {
	code string
	line int
}

func stringencode(str string) string {
	output := strings.Trim(str, " ")
	output = output[1 : len(output)-1]
	return output
}

func makeRegex(str string) *regexp.Regexp {
	Compile, err := regexp.Compile("^(" + str + ")$")
	if err != nil {
		log.Fatal(err)
	}
	return Compile
}

func main() {
	// load .ar file
	path := os.Args[1]
	extention := filepath.Ext(path)
	if extention == "" {
		path += ".ar"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	script := translate(string(data))
	fmt.Println(script)

}
