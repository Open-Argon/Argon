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
var linefunc func(i int, codearray []code) (interface{}, int)

func init() {
	processfunc = translateprocess
	linefunc = translateline
}

func split_by_semicolon_and_newline(str string) []code {
	output := []code{}
	temp := []byte{}
	split := strings.Split(strings.ReplaceAll(str, "\r", ""), "\n")
	splitoutput := []string{}
	for i := 0; i < len(split); i++ {
		if !(skipCompile.MatchString(split[i])) {
			splitoutput = append(splitoutput, split[i])
		}
	}
	str = strings.Join(splitoutput, "\n")
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
	finaloutput := []code{}
	for i := 0; i < len(output); i++ {
		if strings.Trim(output[i].code, " ") != "" {
			finaloutput = append(finaloutput, output[i])
		}
	}
	output = append(output, code{code: string(temp[:]), line: line})
	return output
}

var getCodeInIndent = func(i int, codearray []code) []interface{} {
	var result []interface{}
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
		resp, x := linefunc(i, codearray)
		result = append(result, resp)
		i = x
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
		TYPE := "var"
		if len(firstSplit) > 1 {
			TYPE = firstSplit[0]
		}
		variable := firstSplit[len(firstSplit)-1]
		resp, worked := translateprocess(code{code: VAR[1],
			line: codeseg.line})
		if !worked {
			log.Fatal("invalid value")
		}
		return setVariable{
			TYPE:     TYPE,
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
		first, worked := translateprocess(code{code: split[1],
			line: codeseg.line})
		if worked {
			codeoutput = append([](interface{}){first}, codeoutput...)
		}
		return whileLoop{
			condition: condition,
			code:      codeoutput,
		}, i + len(codeoutput) + 2
	} else if skipCompile.MatchString(codeseg.code) {
	} else {
		err := "Invalid syntax on line "
		log.Fatal("\n\nLine " + fmt.Sprint(codeseg.line+1) + ": " + codeseg.code + "\n" + err + "\n\n")
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
