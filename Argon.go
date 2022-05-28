package main

import (
	"log"
	"os"
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
	importMod(os.Args[1])
}
