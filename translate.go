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
var whileCompile = makeRegex("( *)(while( )+" + anyAndNewline + "+( )+\\[" + anyAndNewline + "*)( *)")
var subCompile = makeRegex("( *)(sub( )+" + variableOnly + "\\(" + anyAndNewline + "*\\)( )+\\[" + anyAndNewline + "*)( *)")
var ifCompile = makeRegex("( *)(if( )+" + anyAndNewline + "+( )+\\[" + anyAndNewline + "*)( *)")
var elseifCompile = makeRegex("( *)(\\]( )+else( )+if( )+" + anyAndNewline + "+( )+\\[" + anyAndNewline + "*)( *)")
var openCompile = makeRegex("( *)((.)+( )+" + anyAndNewline + "+(\\(" + anyAndNewline + "+\\))?( )+\\[" + anyAndNewline + "*)")
var elseCompile = makeRegex("( *)\\]( )+else( )+\\[" + anyAndNewline + "*( *)")
var closeCompile = makeRegex("( *)\\]( *)")
var switchCloseCompile = makeRegex("( *)" + anyAndNewline + "*\\]" + anyAndNewline + "*\\[" + anyAndNewline + "*( *)")
var importCompile = makeRegex("( *)((import " + anyAndNewline + "+)|(from( +)" + anyAndNewline + "+( +)import( +)(( *)" + variableOnly + "( *)(\\,( *)" + variableOnly + "( *))*)))( *)")
var bracketsCompile = makeRegex("( *)\\(" + anyAndNewline + "*\\)( *)")
var functionCompile = makeRegex("( *)" + variableOnly + "\\(" + anyAndNewline + "*\\)( *)")
var variableTextCompile = makeRegex("( *)" + variableOnly + "( *)")
var variableonlyCompile = makeRegex(variableOnly)
var setVariableCompile = makeRegex("( *)(((const|var) (" + variableOnly + "( *))=( *).+)|(( *)" + variableOnly + "( *))(( *)=( *).+))( *)")
var returnStatementCompile = makeRegex("( *)(return(( )+" + anyAndNewline + "*)?)( *)")
var itemsCompiled = makeRegex("( *)(\\[" + anyAndNewline + "*\\])( *)")
var breakStatementCompile = makeRegex("( *)break( *)")
var continueStatementCompile = makeRegex("( *)continue( *)")
var skipCompile = makeRegex("( *)(#(.*))?")
var processfunc func(codeseg code) (any, bool)
var linefunc func(i int, codearray []code) (any, int)
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
		" is greater than or equal to ",
		" is more than or equal to ",
	}, {
		"<",
		" is less than ",
		" is smaller than ",
	}, {
		">",
		" is bigger than ",
		" is greater than ",
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
		" plus ",
		" add ",
		"+",
	}, {
		" minus ",
		" subtract ",
		"-",
	}, {
		"*",
		" x ",
		" times ",
		" multiplied by ",
	}, {
		" mod ",
		" modulo ",
		"%",
	}, {
		" div ",
		" floor division of ",
		" floor divistion ",
		"//",
		"$",
	}, {
		" over ",
		" divided by ",
		"/",
	}, {
		"!**",
		"!^",
		" root ",
		"√",
	}, {
		" to the power of ",
		"^",
		"**",
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
			resp, worked := processfunc(code{code: tempstr, line: line})
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
		resp, worked := processfunc(code{code: tempstr, line: line})
		if worked {
			output = append(output, resp)
		} else {
			return nil, false
		}
	}
	return output, true
}

// make a function that converts backslashes in a string to their correct character
// e.g. "\n" -> newline
// e.g. "\t" -> tab
// e.g. "\r" -> carriage return
// e.g. "\b" -> backspace
// e.g. "\f" -> formfeed
// e.g. "\a" -> alert
// e.g. "\v" -> vertical tab
// e.g. "\0" -> null
// e.g. "\'" -> single quote
// e.g. "\"" -> double quote
// e.g. "\\" -> backslash
// e.g. unicode characters like \u00A0 -> unicode
func unescape(str string) string {
	output := []byte{}
	stringlen := len(str)
	for i := 0; i < stringlen; i++ {
		if str[i] == 92 {
			if i+1 < stringlen {
				i++
				switch str[i] {
				case 110:
					output = append(output, 10)
				case 116:
					output = append(output, 9)
				case 114:
					output = append(output, 13)
				case 98:
					output = append(output, 8)
				case 102:
					output = append(output, 12)
				case 118:
					output = append(output, 11)
				case 120:
					output = append(output, 0)
				case 34:
					output = append(output, 34)
				case 39:
					output = append(output, 39)
				case 92:
					output = append(output, 92)
				case 117:
					if i+4 < stringlen {
						i += 4
						output = append(output, byte(int(str[i])*16*16*16+int(str[i+1])*16*16+int(str[i+2])*16+int(str[i+3])))
					}
				}
			}
		} else {
			output = append(output, str[i])
		}
	}
	return string(output[:])
}

func getParamNames(str string, line int) []string {
	output := []string{}
	temp := []byte{}
	stringlen := len(str)
	for i := 0; i < stringlen; i++ {
		if str[i] == 44 {
			tempstr := strings.Trim(string(temp[:]), " ")
			if !variableonlyCompile.MatchString(tempstr) {
				log.Fatal("invalid param name on line", line+1)
			}
			output = append(output, tempstr)
			temp = []byte{}
		} else {
			temp = append(temp, str[i])
		}
	}
	tempstr := strings.Trim(string(temp[:]), " ")
	if variableonlyCompile.MatchString(tempstr) {
		output = append(output, tempstr)
	}
	return output
}

func split_by_semicolon_and_newline(str string) []code {

	output := []code{}
	temp := []byte{}
	line := 0
	quote := false
	backtick := false
	speech := false
	str = strings.ReplaceAll(str, "\r", "")
	stringlen := len(str)

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
	output = append(output, code{code: string(temp[:]), line: line})
	return output
}

var getCodeInIndent = func(i int, codearray []code, isIf bool) ([]code, int) {
	var result []code
	indent := 0
	startline := codearray[i].line
	result = append(result, code{
		code: strings.SplitN(codearray[i].code, "[", 2)[1],
		line: codearray[i].line,
	})
	i++
	for {
		if i == len(codearray) {
			log.Fatal("invalid opening brackets starting on line ", startline+1)
		}
		if (closeCompile.MatchString(codearray[i].code) && !setVariableCompile.MatchString(codearray[i].code) && !switchCloseCompile.MatchString(codearray[i].code)) || (isIf && (elseCompile.MatchString(codearray[i].code) || elseifCompile.MatchString(codearray[i].code))) {
			indent--
			if indent < 0 {
				fmt.Println(codearray[i].code)
				break
			}
		} else if openCompile.MatchString(codearray[i].code) && !switchCloseCompile.MatchString(codearray[i].code) {
			indent++
		}
		result = append(result, codearray[i])
		i++
	}
	i++
	return result, i
}

var translateprocess = func(codeseg code) (any, bool) {
	for i := 0; i < len(opperators); i++ {
		for j := 0; j < len(opperators[i]); j++ {
			split := strings.Split(codeseg.code, opperators[i][j])
			if len(split) > 1 {
				vals := []any{}
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

						return opperator{
							t:    i,
							vals: vals,
							line: codeseg.line,
						}, true
					} else if len(opperators) == i+1 {
						return nil, false
					}
				}
			}
		}
	}
	if stringCompile.MatchString(codeseg.code) {
		return unescape(stringencode(codeseg.code)), true
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
	} else if itemsCompiled.MatchString(codeseg.code) {
		inside := strings.Trim(codeseg.code, " ")[1 : len(strings.Trim(codeseg.code, " "))-1]
		vals, worked := getValuesFromCommas(inside, codeseg.line)
		if worked {
			return itemsType{
				vals: vals,
				line: codeseg.line,
			}, true
		}
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

var translateline = func(i int, codearray []code) (any, int) {
	codeseg := codearray[i]
	if returnStatementCompile.MatchString(codeseg.code) {
		var val any = strings.Trim(codeseg.code, " ")[6:]
		if val != "" {
			respval, worked := translateprocess(code{code: val.(string), line: codeseg.line})
			if !worked {
				log.Fatal("invalid value on line ", codeseg.line+1)
			}
			val = respval
		} else {
			val = nil
		}
		return returnType{
			val:  val,
			line: codeseg.line,
		}, i + 1
	} else if breakStatementCompile.MatchString(codeseg.code) {
		return breakType{
			line: codeseg.line,
		}, i + 1
	} else if continueStatementCompile.MatchString(codeseg.code) {
		return continueType{
			line: codeseg.line,
		}, i + 1
	} else if skipCompile.MatchString(codeseg.code) {
		return nil, i + 1
	}
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
			log.Fatal("invalid value on line ", codeseg.line+1, ": ", VAR[1])
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
			log.Fatal("invalid value on line ", codeseg.line+1)
		}
		codedata, x := getCodeInIndent(i, codearray, false)
		codeoutput := []any{}

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
			log.Fatal("invalid value on line ", codeseg.line+1, ": ", split[0])
		}
		codedata, x := getCodeInIndent(i, codearray, true)
		codeoutput := []any{}

		codelen := len(codedata)
		for i := 0; i < codelen; {
			resp, x := linefunc(i, codedata)
			i = x
			codeoutput = append(codeoutput, resp)
		}
		ifstatementcode := []iftype{
			{
				condition: condition,
				code:      codeoutput,
			},
		}
		elsecode := []any{}
		for {
			if elseifCompile.MatchString(codearray[x-1].code) {
				str := strings.Trim(codearray[x-1].code, " ")
				str = strings.SplitN(str[1:len(codearray[x-1].code)-1], "if", 2)[1]
				condition, worked = translateprocess(code{code: str,
					line: codearray[x-1].line})
				if !worked {
					log.Fatal("invalid value on line ", codearray[x-1].line+1, ": ", split[0])
				}
				codedata, j := getCodeInIndent(x-1, codearray, true)
				x += j - 5
				codeoutput := []any{}

				codelen := len(codedata)
				for i := 0; i < codelen; {
					resp, x := linefunc(i, codedata)
					i = x
					codeoutput = append(codeoutput, resp)
				}
				ifstatementcode = append(ifstatementcode, iftype{
					condition: condition,
					code:      codeoutput,
				})
			} else if elseCompile.MatchString(codearray[x-1].code) {
				codedata, j := getCodeInIndent(x-1, codearray, false)
				x += j - 3
				codelen := len(codedata)
				for i := 0; i < codelen; {
					resp, x := linefunc(i, codedata)
					i = x
					elsecode = append(elsecode, resp)
				}
				break
			} else {
				break
			}
		}
		fmt.Println(ifstatementcode)
		fmt.Println(elsecode)
		return ifstatement{
			statments: ifstatementcode,
			FALSE:     elsecode,
		}, x
	} else if subCompile.MatchString(codeseg.code) {
		sub := strings.Trim(codeseg.code, " ")
		sub = sub[4:]
		split := strings.SplitN(sub, "[", 2)
		functioninfo := strings.SplitN(strings.Trim(split[0], " "), "(", 2)
		functionname := functioninfo[0]
		argstr := functioninfo[1]
		argstr = argstr[:len(argstr)-1]
		args := getParamNames(argstr, codeseg.line)
		codedata, x := getCodeInIndent(i, codearray, true)
		codeoutput := []any{}

		codelen := len(codedata)
		for i := 0; i < codelen; {
			resp, x := linefunc(i, codedata)
			i = x
			codeoutput = append(codeoutput, resp)
		}
		return setFunction{
			name: functionname,
			args: args,
			code: codeoutput,
		}, x
	} else if importCompile.MatchString(codeseg.code) {
		str := strings.Trim(codeseg.code, " ")
		var toImport any = nil
		var path any
		if strings.HasPrefix(str, "import") {
			p, worked := processfunc(code{code: str[7:], line: codeseg.line})
			path = p
			if !worked {
				log.Fatal("invalid import path on line ", codeseg.line+1)
			}
		} else {
			info := strings.SplitN(str[5:], " import ", 2)
			p, worked := processfunc(code{code: info[0], line: codeseg.line})
			path = p
			if !worked {
				log.Fatal("invalid import path on line ", codeseg.line+1)
			}
			toImport = getParamNames(info[1], codeseg.line)
		}
		return importType{
			path:     path,
			toImport: toImport,
			line:     codeseg.line,
		}, i + 1
	} else {
		err := "Invalid syntax on line "
		log.Fatal("\n\nLine " + fmt.Sprint(codeseg.line+1) + ": " + codeseg.code + "\n" + err + "\n\n")
	}
	return nil, i + 1
}

var translate = func(str string) []any {
	code := split_by_semicolon_and_newline(str)
	codelen := len(code)
	output := []any{}
	i := 0
	for i < codelen {
		resp, x := translateline(i, code)
		i = x
		output = append(output, resp)
	}
	return output
}
