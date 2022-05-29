package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
)

var runop func(codeseg any) any

func init() {
	runop = runprocess
}

func run(lines []any) {
	for i := 0; i < len(lines); i++ {
		if lines[i] != nil {
			runprocess(lines[i])
		}
	}
}

func anyToArgon(x any) string {
	switch x.(type) {
	case string:
		return x.(string)
	case float64:
		return fmt.Sprint(x)
	case bool:
		if x.(bool) {
			return "yes"
		} else {
			return "no"
		}
	case nil:
		return "unknown"
	default:
		return fmt.Sprint(x)
	}
}

func runprocess(codeseg any) any {
	switch codeseg.(type) {
	case opperator:
		return runOperator(codeseg.(opperator))
	case variable:
		myvar := vars[codeseg.(variable).variable]
		if myvar.EXISTS != nil {
			return myvar.VAL
		}
		log.Fatal("undecared variable " + codeseg.(variable).variable + " on line " + fmt.Sprint(codeseg.(variable).line+1))
	case funcCallType:
		return callFunc(codeseg.(funcCallType))
	case whileLoop:
		whileloop := codeseg.(whileLoop)
		for boolean(runop(whileloop.condition)) {
			run(whileloop.code)
		}
	case ifstatement:
		iff := codeseg.(ifstatement)
		if boolean(runop(iff.condition)) {
			run(iff.TRUE)
		} else {
			run(iff.FALSE)
		}
	case importType:
		importMod((codeseg.(importType).path).(string))
	case setVariable:
		setVariableVal(codeseg.(setVariable))
	}
	return codeseg
}

func callFunc(call funcCallType) any {
	if vars[call.name].EXISTS == nil {
		log.Fatal("undecared function on line " + fmt.Sprint(call.line+1))
	} else if vars[call.name].TYPE != "func" && vars[call.name].TYPE != "init" {
		log.Fatal("cannot call " + vars[call.name].TYPE + " on line " + fmt.Sprint(call.line+1))
	} else if vars[call.name].FUNC == false {
		log.Fatal(call.name + " is not a function on line " + fmt.Sprint(call.line+1))
	}
	argvals := []any{}
	for i := 0; i < len(call.args); i++ {
		argvals = append(argvals, runprocess(call.args[i]))
	}
	return vars[call.name].VAL.(func(...any) any)(argvals...)
}

func setVariableVal(x setVariable) {
	variable := vars[x.variable]
	if variable.EXISTS == nil || variable.TYPE == "var" {
		vars[x.variable] = variableValue{
			TYPE:   x.TYPE,
			EXISTS: true,
			VAL:    runop(x.value),
			FUNC:   false,
		}
	} else {
		log.Fatal("cannot edit " + variable.TYPE + " variable on line " + fmt.Sprint(x.line+1))
	}
}

func dynamicAdd(x any, y any) any {
	stringconvert := false
	switch x.(type) {
	case string:
		stringconvert = true
	}
	switch y.(type) {
	case string:
		stringconvert = true
	}
	if stringconvert {
		return fmt.Sprint(x) + fmt.Sprint(y)
	}
	xnum := number(x)
	ynum := number(y)
	return xnum + ynum
}

func xiny(x any, y []any) bool {
	for i := 0; i < len(y); i++ {
		if x == y[i] {
			return true
		}
	}
	return false
}

func boolean(x any) bool {
	return (x != false && x != nil && x != 0 && x != "")
}

func number(x any) float64 {
	switch x.(type) {
	case float64:
		return x.(float64)
	case float32:
		return float64(x.(float32))
	case int:
		return float64(x.(int))
	case bool:
		if x.(bool) {
			return 1
		} else {
			return 0
		}
	default:
		num, err := strconv.ParseFloat(fmt.Sprint(x), 64)
		if err != nil {
			log.Fatal(err)
		}
		return num
	}
}

func runOperator(opperation opperator) any {
	var output any
	for i := 0; i < len(opperation.vals); i++ {
		switch opperation.t {
		case 0:
			x := runop(opperation.vals[i])
			if output == nil {
				output = x
			} else {
				if boolean(output) && boolean(x) {
					output = x
				} else {
					output = false
					break
				}
			}
		case 1:
			x := runop(opperation.vals[i])
			if boolean(x) {
				output = x
				break
			}
		case 2:
			x := runop(opperation.vals[i])
			if output == nil {
				output = x
			} else {
				output = xiny(output, x.([]any))
			}
		case 3:
			x := runop(opperation.vals[i])
			if output == nil {
				output = x
			} else {
				output = !xiny(output, x.([]interface{}))
			}
		case 4:
			x := runop(opperation.vals[i])
			if output == nil {
				output = x
			} else {
				output = (number(output) <= number(x))
			}
		case 5:
			x := runop(opperation.vals[i])
			if output == nil {
				output = x
			} else {
				output = (number(output) >= number(x))
			}
		case 6:
			x := runop(opperation.vals[i])
			if output == nil {
				output = x
			} else {
				output = (number(output) < number(x))
			}
		case 7:
			x := runop(opperation.vals[i])
			if output == nil {
				output = x
			} else {
				output = (number(output) > number(x))
			}
		case 8:
			x := runop(opperation.vals[i])
			if output == nil {
				output = x
			} else {
				output = (output != x)
			}
		case 9:
			x := runop(opperation.vals[i])
			if output == nil {
				output = x
			} else {
				output = (output == x)
			}
		case 10:
			x := runop(opperation.vals[i])
			if output == nil {
				output = x
			} else {
				output = (number(output) - number(x))
			}
		case 11:
			x := runop(opperation.vals[i])
			if output == nil {
				output = x
			} else {
				output = dynamicAdd(output, x)
			}
		case 12:
			x := runop(opperation.vals[i])
			if output == nil {
				output = x
			} else {
				output = math.Pow(number(output), 1/number(x))
			}
		case 13:
			x := runop(opperation.vals[i])
			if output == nil {
				output = x
			} else {
				output = math.Pow(number(output), number(x))
			}
		case 14:
			x := runop(opperation.vals[i])
			if output == nil {
				output = x
			} else {
				output = (number(output) * number(x))
			}
		case 15:
			x := runop(opperation.vals[i])
			if output == nil {
				output = x
			} else {
				output = math.Floor(number(output) / number(x))
			}
		case 16:
			x := runop(opperation.vals[i])
			if output == nil {
				output = x
			} else {
				output = math.Mod(number(output), number(x))
			}
		case 17:
			x := runop(opperation.vals[i])
			if output == nil {
				output = x
			} else {
				output = (number(output) / number(x))
			}
		}
	}
	return output
}
