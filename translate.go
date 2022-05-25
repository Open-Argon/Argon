package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

var stringCompile = makeRegex("( *)((\"((\\\\([a-z\\\"']))|[^\\\"])*\"))|(('((\\\\([a-z\\'\"]))|[^\\'])*'))( *)")
var numberCompile = makeRegex("( *)(\\-)?(([0-9]+(\\.[0-9]+)?)(e((\\-|\\+)?([0-9]+(\\.[0-9]+)?)))?)( *)")
var variableOnly = "([a-zA-Z])([a-zA-Z0-9])*"
var variableTextCompile = makeRegex("( *)" + variableOnly + "( *)")
var setVariableCompile = makeRegex("( *)(((const|var) ({" + fmt.Sprint(variableTextCompile) + "})(( +)=( +).+)?)|(" + fmt.Sprint(variableTextCompile) + ")(( +)=( +).+))( *)")
var skipCompile = makeRegex("( *)(#( *))?")
var translateprocess = func(codeseg code) (interface{}, bool) {
	if stringCompile.MatchString(codeseg.code) {
		return (stringencode(codeseg.code)), true
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
var translateline = func(i int, code []code) (interface{}, int) {
	codeseg := code[i]
	resp, worked := translateprocess(codeseg)
	if worked {
		return resp, i + 1
	} else if setVariableCompile.MatchString(codeseg.code) {

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
		output = append(output, resp)
	}
	return output
}
