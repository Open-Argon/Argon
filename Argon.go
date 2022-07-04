package main

import (
	"log"
	"os"
	"regexp"
	"strings"
	"time"
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
	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	ArgonSetSeed(time.Now().UnixMilli())
	if len(os.Args) == 1 {
		shell()
	} else {
		_, err := importMod(os.Args[1], ex)
		if err != nil {
			log.Fatal(err)
		}
	}
}
