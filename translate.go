package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

var stringCompile = makeRegex("( *)((\"((\\\\([a-z\\\"']))|[^\\\"])*\"))|(('((\\\\([a-z\\'\"]))|[^\\'])*'))( *)")
var numberCompile = makeRegex("( *)(\\-)?(([0-9]+(\\.[0-9]+)?)(e((\\-|\\+)?([0-9]+(\\.[0-9]+)?)))?)( *)")
var whileCompile = makeRegex("( *)(while .+ \\[.*)( *)")
var openCompile = makeRegex("( *)([a-z]+ .+ \\[.*)( *)")
var closeCompile = makeRegex("( *)\\]( *)")
var bracketsCompile = makeRegex("( *)\\(.*\\)( *)")
var variableOnly = "([a-zA-Z])([a-zA-Z0-9])*"
var variableTextCompile = makeRegex("( *)" + variableOnly + "( *)")
var setVariableCompile = makeRegex("( *)(((const|var) (" + variableOnly + "( *))=( *).+)|(( *)" + variableOnly + "( *))(( *)=( *).+))( *)")
var skipCompile = makeRegex("( *)(#( *))?")
var processfunc func(codeseg code) (interface{}, bool)

func init() {
	processfunc = translateprocess
}

var getCodeInIndent = func(i int, codearray []code) []code {
	var result = []code{}
	indent := 0
	for true {
		if closeCompile.MatchString(codearray[i].code) {
			indent++
			if indent >= 0 {
				break
			}
		} else if openCompile.MatchString(codearray[i].code) {
			indent--
		}
		result = append(result, codearray[i])
		i++
	}
	return result
}

var translateprocess = func(codeseg code) (interface{}, bool) {
	if stringCompile.MatchString(codeseg.code) {
		return (stringencode(codeseg.code)), true
	} else if bracketsCompile.MatchString(codeseg.code) {
		shorter := strings.Trim(codeseg.code, " ")
		shorter = shorter[1 : len(shorter)-1]
		brackets, ran := processfunc(code{
			code: shorter,
			line: codeseg.line})
		return brackets, ran
	} else if numberCompile.MatchString(codeseg.code) {
		if s, err := strconv.ParseFloat(strings.Trim(codeseg.code, " "), 64); err == nil {
			return s, true
		} else {
			log.Fatal(err)
		}
	} else if variableTextCompile.MatchString(codeseg.code) {
		return variable{
			variable: strings.Trim(codeseg.code, " "),
			line:     codeseg.line,
		}, true
	}
	return nil, false
}

var translateline = func(i int, codearray []code) (interface{}, int) {
	codeseg := codearray[i]
	resp, worked := translateprocess(codeseg)
	if worked {
		return resp, i + 1
	} else if setVariableCompile.MatchString(codeseg.code) {
		VAR := strings.Split(codeseg.code, "=")
		firstSplit := strings.Split(strings.Trim(VAR[0], " "), " ")
		variable := firstSplit[len(firstSplit)-1]
		resp, worked := translateprocess(code{code: VAR[1],
			line: codeseg.line})
		if !worked {
			log.Fatal("invalid value")
		}
		return setVariable{
			variable: variable,
			value:    resp,
			line:     codeseg.line,
		}, i + 1
	} else if whileCompile.MatchString(codeseg.code) {
		str := strings.Trim(codeseg.code, " ")
		strlen := len(str)
		str = str[5:strlen]
		split := strings.SplitN(str, "[", 2)
		condition, worked := translateprocess(code{code: split[0],
			line: codeseg.line})
		if !worked {
			log.Fatal("invalid value on line", codeseg.line)
		}
		codeoutput := getCodeInIndent(i+1, codearray)
		codeoutput = append([]code{
			{
				code: split[1],
				line: codeseg.line,
			},
		}, codeoutput...)
		return whileLoop{
			condition: condition,
			code:      codeoutput,
		}, i + 1 + len(codeoutput)
	} else if skipCompile.MatchString(codeseg.code) {
	} else {
		err := "Invalid syntax on line " + fmt.Sprint(codeseg.line+1)
		log.Fatal(err)
	}
	return nil, i + 1
}

var translate = func(str string) []interface{} {
	code := split_by_semicolon_and_newline(str)
	codelen := len(code)
	output := [](interface{}){}
	for i := 0; i < codelen; {
		resp, x := translateline(i, code)
		i = x
		if resp != nil {
			output = append(output, resp)
		}
	}
	return output
}
