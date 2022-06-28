package main

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
)

var runop func(codeseg any, localvars map[string]variableValue) (any, any)

func init() {
	runop = runprocess
}

func run(lines []any, localvars map[string]variableValue) (any, any, [][]any) {
	output := ([][]any{})
	for i := 0; i < len(lines); i++ {
		if lines[i] != nil {
			val, ty := runprocess(lines[i], localvars)
			output = append(output, []any{val, ty})
			if ty == "return" || ty == "break" {
				return ty, val, output
			}
		}
	}
	return nil, nil, output
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
	case variable:
		return x.(variable).variable
	default:
		return fmt.Sprint(x)
	}
}

func runprocess(codeseg any, localvars map[string]variableValue) (any, any) {
	switch codeseg.(type) {
	case opperator:
		return runOperator(codeseg.(opperator), localvars), "opperator"
	case variable:
		myvar := localvars[codeseg.(variable).variable]
		if myvar.EXISTS == nil {
			myvar = vars[codeseg.(variable).variable]
		}
		if myvar.EXISTS != nil {
			if myvar.TYPE != "func" && myvar.TYPE != "init_function" {
				return myvar.VAL, "value"
			} else {
				return codeseg, "value"
			}
		}
		log.Fatal("undecared variable " + codeseg.(variable).variable + " on line " + fmt.Sprint(codeseg.(variable).line+1))
	case funcCallType:
		return callFunc(codeseg.(funcCallType), localvars), "value"
	case whileLoop:
		whileloop := codeseg.(whileLoop)
		resp, _ := runop(whileloop.condition, localvars)
		for boolean(resp) {
			ty, val, _ := run(whileloop.code, localvars)
			if ty == "break" {
				break
			} else if ty == "return" {
				return val, "return"
			}
		}
	case ifstatement:
		iff := codeseg.(ifstatement)
		resp, _ := runop(iff.condition, localvars)
		if boolean(resp) {
			ty, val, _ := run(iff.TRUE, localvars)
			return val, ty
		} else {
			ty, val, _ := run(iff.FALSE, localvars)
			return val, ty
		}
	case importType:
		importMod((codeseg.(importType).path).(string))
		return nil, nil
	case setVariable:
		setVariableVal(codeseg.(setVariable), localvars)
		return nil, nil
	case setFunction:
		setFunctionVal(codeseg.(setFunction), localvars)
		return nil, nil
	case returnType:
		val, worked := runop(codeseg.(returnType).val, localvars)
		if worked == nil {
			log.Fatal("return statement must return a value on line " + fmt.Sprint(codeseg.(returnType).line+1))
		}
		return val, "return"
	case breakType:
		return nil, "break"
	}
	return codeseg, "value"
}

func callFunc(call funcCallType, localvars map[string]variableValue) any {
	variables := localvars
	if variables[call.name].EXISTS == nil {
		variables = vars
	}
	if variables[call.name].EXISTS == nil {
		log.Fatal("undecared function '" + call.name + "' on line " + fmt.Sprint(call.line+1))
	}
	callable := call
	for variables[callable.name].TYPE != "func" && vars[callable.name].TYPE != "init_function" && fmt.Sprint(reflect.TypeOf(variables[callable.name].VAL)) == "main.variable" {
		callable = funcCallType{name: variables[callable.name].VAL.(variable).variable, args: callable.args, line: callable.line}
	}
	if variables[callable.name].TYPE != "func" && variables[callable.name].TYPE != "init_function" {
		log.Fatal("'" + call.name + "' is not a function on line " + fmt.Sprint(call.line+1))
	}
	if vars[callable.name].FUNC {

		argvals := []any{}
		for i := 0; i < len(callable.args); i++ {
			resp, _ := runprocess(callable.args[i], localvars)
			argvals = append(argvals, resp)
		}
		return vars[callable.name].VAL.(func(...any) any)(argvals...)
	} else {
		argvars := make(map[string]variableValue)
		for i := 0; i < len(callable.args); i++ {
			resp, _ := runprocess(callable.args[i], localvars)
			name := variables[callable.name].VAL.(setFunction).args[i]
			argvars[name] = variableValue{
				VAL:    resp,
				TYPE:   "var",
				EXISTS: true,
				FUNC:   false,
			}
		}
		ty, val, _ := run(vars[callable.name].VAL.(setFunction).code, argvars)
		if ty != "return" && ty != nil {
			log.Fatal(ty, " is not allowed in top level of function on line ", (callable.line + 1))
		}
		return val
	}
}

func setVariableVal(x setVariable, localvars map[string]variableValue) {
	variable := localvars
	if variable[x.variable].EXISTS == nil {
		variable = vars
	}
	vari := variable[x.variable]
	if variable[x.variable].EXISTS == nil {
		resp, _ := runop(x.value, localvars)
		variable[x.variable] = variableValue{
			TYPE:   x.TYPE,
			EXISTS: true,
			VAL:    resp,
			FUNC:   false,
		}
	} else if variable[x.variable].TYPE == "var" {
		val, _ := runop(x.value, localvars)
		variable[x.variable] = variableValue{
			TYPE:   x.TYPE,
			EXISTS: variable[x.variable].EXISTS,
			VAL:    val,
			FUNC:   variable[x.variable].FUNC,
		}
	} else {
		log.Fatal("cannot edit " + vari.TYPE + " variable on line " + fmt.Sprint(x.line+1))
	}
}

func setFunctionVal(x setFunction, localvars map[string]variableValue) {
	variable := localvars[x.name]
	if variable.EXISTS == nil {
		variable = vars[x.name]
	}
	if vars[x.name].EXISTS == nil {
		vars[x.name] = variableValue{
			TYPE:   "func",
			EXISTS: true,
			VAL:    x,
			FUNC:   false,
		}
	} else if vars[x.name].TYPE == "func" {
		variable.VAL = x
		variable.FUNC = false
	} else {
		log.Fatal("cannot edit " + vars[x.name].TYPE + " variable on line " + fmt.Sprint(x.line+1))
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

func runOperator(opperation opperator, localvars map[string]variableValue) any {
	var output any
	for i := 0; i < len(opperation.vals); i++ {
		switch opperation.t {
		case 0:
			x, _ := runop(opperation.vals[i], localvars)
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
			x, _ := runop(opperation.vals[i], localvars)
			if boolean(x) {
				output = x
				break
			}
		case 2:
			x, _ := runop(opperation.vals[i], localvars)
			if output == nil {
				output = x
			} else {
				output = xiny(output, x.([]any))
			}
		case 3:
			x, _ := runop(opperation.vals[i], localvars)
			if output == nil {
				output = x
			} else {
				output = !xiny(output, x.([]interface{}))
			}
		case 4:
			x, _ := runop(opperation.vals[i], localvars)
			if output == nil {
				output = x
			} else {
				output = (number(output) <= number(x))
			}
		case 5:
			x, _ := runop(opperation.vals[i], localvars)
			if output == nil {
				output = x
			} else {
				output = (number(output) >= number(x))
			}
		case 6:
			x, _ := runop(opperation.vals[i], localvars)
			if output == nil {
				output = x
			} else {
				output = (number(output) < number(x))
			}
		case 7:
			x, _ := runop(opperation.vals[i], localvars)
			if output == nil {
				output = x
			} else {
				output = (number(output) > number(x))
			}
		case 8:
			x, _ := runop(opperation.vals[i], localvars)
			if output == nil {
				output = x
			} else {
				output = (output != x)
			}
		case 9:
			x, _ := runop(opperation.vals[i], localvars)
			if output == nil {
				output = x
			} else {
				output = (output == x)
			}
		case 10:
			x, _ := runop(opperation.vals[i], localvars)
			if output == nil {
				output = x
			} else {
				output = (number(output) - number(x))
			}
		case 11:
			x, _ := runop(opperation.vals[i], localvars)
			if output == nil {
				output = x
			} else {
				output = dynamicAdd(output, x)
			}
		case 12:
			x, _ := runop(opperation.vals[i], localvars)
			if output == nil {
				output = x
			} else {
				output = (number(output) * number(x))
			}
		case 13:
			x, _ := runop(opperation.vals[i], localvars)
			if output == nil {
				output = x
			} else {
				output = math.Floor(number(output) / number(x))
			}
		case 14:
			x, _ := runop(opperation.vals[i], localvars)
			if output == nil {
				output = x
			} else {
				output = math.Mod(number(output), number(x))
			}
		case 15:
			x, _ := runop(opperation.vals[i], localvars)
			if output == nil {
				output = x
			} else {
				output = (number(output) / number(x))
			}
		case 16:
			x, _ := runop(opperation.vals[i], localvars)
			if output == nil {
				output = x
			} else {
				output = math.Pow(number(output), 1/number(x))
			}
		case 17:
			x, _ := runop(opperation.vals[i], localvars)
			if output == nil {
				output = x
			} else {
				output = math.Pow(number(output), number(x))
			}
		}
	}
	return output
}
