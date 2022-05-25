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

func split_by_semicolon_and_newline(str string) []code {
	output := []code{}
	temp := []byte{}
	stringlen := len(str)
	line := 0
	for i := 0; i < stringlen; i++ {
		if str[i] == 59 || str[i] == 10 {
			output = append(output, code{code: string(temp[:]), line: line})
			if str[i] == 10 {
				line++
			}
			temp = []byte{}
		} else {
			temp = append(temp, str[i])
		}
	}
	output = append(output, code{code: string(temp[:]), line: line})
	return output
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
