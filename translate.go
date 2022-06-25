package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

var anyAndNewline = "((.)|(\\n))"
var variableOnly = "([a-zA-Z_])([a-zA-Z0-9_])*"
var stringCompile = makeRegex("(( *)\"((\\\\([a-z\\\"'`]))|[^\\\"])*\"( *))|(( *)'((\\\\([a-z\\'\"`]))|[^\\'])*'( *))|(( *)(`(\\\\[a-z\\\"'`\\n]|[^\\`])*`)( *))")
var numberCompile = makeRegex("( *)(\\-)?(([0-9]+(\\.[0-9]+)?)(e((\\-|\\+)?([0-9]+(\\.[0-9]+)?)))?)( *)")
var whileCompile = makeRegex("( *)(while " + anyAndNewline + "+ \\[" + anyAndNewline + "*)( *)")
var subCompile = makeRegex("( *)(sub " + variableOnly + "\\(" + anyAndNewline + "*\\) \\[" + anyAndNewline + "*)( *)")
var ifCompile = makeRegex("( *)(if " + anyAndNewline + "+ \\[" + anyAndNewline + "*)( *)")
var elseifCompile = makeRegex("( *)(\\] else if " + anyAndNewline + "+ \\[" + anyAndNewline + "*)( *)")
var openCompile = makeRegex("( *)([a-z]+ " + anyAndNewline + "+ \\[.*)( *)")
var elseCompile = makeRegex("( *)\\] else \\[" + anyAndNewline + "*( *)")
var closeCompile = makeRegex("( *)\\]( *)")
var importCompile = makeRegex("( *)import " + anyAndNewline + "+( *)")
var bracketsCompile = makeRegex("( *)\\(" + anyAndNewline + "*\\)( *)")
var functionCompile = makeRegex("( *)" + variableOnly + "\\(" + anyAndNewline + "*\\)( *)")
var variableTextCompile = makeRegex("( *)" + variableOnly + "( *)")
var setVariableCompile = makeRegex("( *)(((const|var) (" + variableOnly + "( *))=( *).+)|(( *)" + variableOnly + "( *))(( *)=( *).+))( *)")
var skipCompile = makeRegex("( *)(#( *))?")
var processfunc func(codeseg code) (interface{}, bool)
var linefunc func(i int, codearray []code) (interface{}, int)
var opperators = [][]string{
	{
		"&&",
		" and ",
	}, {
		"||",
		" or ",
	}, {
		"!@",
		" not in ",
		" is not in ",
		" isnt in ",
	}, {
		"@",
		" in ",
		" is in ",
	}, {
		"<=",
		" is less than or equal to ",
		" is smaller than or equal to ",
	}, {
		">=",
		" is bigger than or equal to ",
		" is more than or equal to ",
	}, {
		"<",
		" is less than ",
		" is smaller than ",
	}, {
		">",
		" is bigger than ",
		" is more than ",
	}, {
		"!=",
		" is not ",
		" is not equal to ",
		" isnt equal to ",
		" isnt ",
		"!==",
	}, {
		"==",
		" equals ",
		" is equal to ",
		" is ",
		"===",
	}, {
		" minus ",
		" subtract ",
		"-",
	}, {
		" plus ",
		" add ",
		"+",
	}, {
		"!**",
		"!^",
		" root ",
		"√",
	}, {
		" to the power of ",
		"^",
		"**",
	}, {
		"*",
		" x ",
		" times ",
		" multiplied by ",
	}, {
		" div ",
		" floor division of ",
		" floor divistion ",
		"//",
		"$",
	}, {
		" mod ",
		" modulo ",
		"%",
	}, {
		" over ",
		" divided by ",
		"/",
	}}

func init() {
	processfunc = translateprocess
	linefunc = translateline
}

func getValuesFromCommas(str string, line int) ([]any, bool) {
	output := []any{}
	temp := []byte{}
	stringlen := len(str)
	for i := 0; i < stringlen; i++ {
		if str[i] == 44 {
			tempstr := string(temp[:])
			resp, worked := processfunc(code{code: tempstr, line: 0})
			if worked {
				output = append(output, resp)
				temp = []byte{}
				continue
			}
		}
		temp = append(temp, str[i])
	}
	tempstr := string(temp[:])
	if tempstr != "" {
		resp, worked := processfunc(code{code: tempstr, line: 0})
		if worked {
			output = append(output, resp)
		} else {
			return nil, false
		}
	}
	return output, true
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
	quote := false
	backtick := false
	speech := false

	for i := 0; i < stringlen; i++ {
		if str[i] == 34 && (len(temp) == 0 || str[i-1] != 92) && !speech && !backtick {
			quote = !quote
		} else if str[i] == 96 && (len(temp) == 0 || str[i-1] != 92) && !speech && !quote {
			backtick = !backtick
		} else if str[i] == 39 && (len(temp) == 0 || str[i-1] != 92) && !quote && !backtick {
			speech = !speech
		}
		if (str[i] == 59 && !quote && !speech && !backtick) || (str[i] == 10 && !backtick) {
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

var getCodeInIndent = func(i int, codearray []code, isIf bool) ([]code, int) {
	var result []code
	indent := 0
	result = append(result, code{
		code: strings.SplitN(codearray[i].code, "[", 2)[1],
		line: codearray[i].line,
	})
	i++
	for {
		if closeCompile.MatchString(codearray[i].code) || (isIf && (elseCompile.MatchString(codearray[i].code) || elseifCompile.MatchString(codearray[i].code))) {
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
	i++
	return result, i
}

var translateprocess = func(codeseg code) (interface{}, bool) {
	for i := 0; i < len(opperators); i++ {
		for j := 0; j < len(opperators[i]); j++ {
			split := strings.Split(codeseg.code, opperators[i][j])
			if len(split) > 1 {
				vals := []interface{}{}
				current := 0
				for k := 0; k < len(split); k++ {
					currentstr := strings.Join(split[current:k], opperators[i][j])
					val, worked := processfunc(code{
						code: currentstr,
						line: codeseg.line,
					})
					if worked {
						vals = append(vals, val)
						current = k
					}
				}
				if len(vals) >= 1 {
					currentstr := strings.Join(split[current:], opperators[i][j])
					val, worked := processfunc(code{
						code: currentstr,
						line: codeseg.line,
					})
					if worked {
						vals = append(vals, val)
					}
					return opperator{
						t:    i,
						vals: vals,
						line: codeseg.line,
					}, true
				}
			}
		}
	}
	if stringCompile.MatchString(codeseg.code) {
		return (stringencode(codeseg.code)), true
	} else if functionCompile.MatchString(codeseg.code) {
		str := strings.Trim(codeseg.code, " ")
		bracketsplit := strings.SplitN(str, "(", 2)
		funcname := strings.Trim(bracketsplit[0], " ")
		params := bracketsplit[1][:len(bracketsplit[1])-1]
		paramstranslated, worked := getValuesFromCommas(params, codeseg.line)
		if worked {
			return funcCallType{
				name: funcname,
				args: paramstranslated,
				line: codeseg.line,
			}, true
		}
	} else if bracketsCompile.MatchString(codeseg.code) {
		shorter := strings.Trim(codeseg.code, " ")
		shorter = shorter[1 : len(shorter)-1]
		brackets, ran := processfunc(code{
			code: shorter,
			line: codeseg.line,
		})
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
			log.Fatal("invalid value on line", codeseg.line+1)
		}
		codedata, x := getCodeInIndent(i, codearray, false)
		codeoutput := []interface{}{}

		codelen := len(codedata)
		for i := 0; i < codelen; {
			resp, x := linefunc(i, codedata)
			i = x
			codeoutput = append(codeoutput, resp)
		}

		return whileLoop{
			condition: condition,
			code:      codeoutput,
		}, x
	} else if ifCompile.MatchString(codeseg.code) {
		str := strings.Trim(codeseg.code, " ")
		strlen := len(str)
		str = str[2:strlen]
		split := strings.SplitN(str, "[", 2)
		condition, worked := translateprocess(code{code: split[0],
			line: codeseg.line})
		if !worked {
			log.Fatal("invalid value on line", codeseg.line+1)
		}
		codedata, x := getCodeInIndent(i, codearray, true)
		codeoutput := []interface{}{}

		codelen := len(codedata)
		for i := 0; i < codelen; {
			resp, x := linefunc(i, codedata)
			i = x
			codeoutput = append(codeoutput, resp)
		}
		elsecode := []any{}
		elsecodelen := 0
		if elseCompile.MatchString(codearray[x-1].code) {
			codedata, x := getCodeInIndent(x-1, codearray, false)
			elsecodelen = x - 3
			codelen := len(codedata)
			for i := 0; i < codelen; {
				resp, x := linefunc(i, codedata)
				i = x
				elsecode = append(elsecode, resp)
			}
		}

		return ifstatement{
			condition: condition,
			TRUE:      codeoutput,
			FALSE:     elsecode,
		}, x + elsecodelen
	} else if subCompile.MatchString(codeseg.code) {
		sub := strings.Trim(codeseg.code, " ")
		sub = sub[3:]
		sub = strings.Trim(sub, " ")
		fmt.Println(sub)
		return ifstatement{
			condition: []interface{}{},
			TRUE:      []interface{}{},
			FALSE:     []interface{}{},
		}, i + 1
	} else if importCompile.MatchString(codeseg.code) {
		str := strings.Trim(codeseg.code, " ")
		path, worked := processfunc(code{code: str[7:], line: codeseg.line})
		if !worked {
			log.Fatal("invalid import path on line ", codeseg.line+1)
		}
		return importType{
			path: path,
			line: codeseg.line,
		}, i + 1
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
	output := []interface{}{}
	i := 0
	for i < codelen {
		resp, x := translateline(i, code)
		i = x
		output = append(output, resp)
	}
	return output
}
