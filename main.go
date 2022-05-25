package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type process struct {
	t    int
	x    interface{}
	y    interface{}
	line int
}

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

func main() {
	stringCompile, err := regexp.Compile("^(( *)((\"((\\\\([a-z\\\"']))|[^\\\"])*\"))|(('((\\\\([a-z\\'\"]))|[^\\'])*'))( *))$")
	if err != nil {
		log.Fatal(err)
	}
	numberCompile, err := regexp.Compile("^(( *)(\\-)?(([0-9]+(\\.[0-9]+)?)(e(([0-9]+(\\.[0-9]+)?)))?)( *))$")
	if err != nil {
		log.Fatal(err)
	}
	translateline := func(codeseg code) interface{} {
			if stringCompile.MatchString(codeseg.code) {
				return (stringencode(codeseg.code))
			} else if numberCompile.MatchString(codeseg.code) {
				if s, err := strconv.ParseFloat(strings.Trim(codeseg.code, " "), 64); err == nil {
					return (s)
				} else {
					log.Fatal(err)
				}
			} else {
				err := "Invalid syntax on line " + fmt.Sprint(codeseg.line+1)
				log.Fatal(err)
			}
      return nil
		}

  translate := func(str string) []interface{} {
		code := split_by_semicolon_and_newline(str)
		codelen := len(code)
    output := [](interface{}){}
		for i := 0; i < codelen; i++ {
      output = append(output, translateline(code[i]))
    }
    return output
  }
	script := translate(" 10 ")
  fmt.Println(script)
}
