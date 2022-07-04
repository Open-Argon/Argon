package main

import (
	"fmt"
	"log"
	"math"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

var runop func(codeseg any, origin string, vargroups []map[string]variableValue) (any, any)

func init() {
	runop = runprocess
}

func run(lines []any, origin string, vargroups []map[string]variableValue) (any, any, [][]any) {
	output := ([][]any{})
	for i := 0; i < len(lines); i++ {
		if lines[i] != nil {
			val, ty := runprocess(lines[i], origin, vargroups)
			output = append(output, []any{val, ty})
			if ty == "return" || ty == "break" || ty == "continue" {
				return ty, val, output
			}
		}
	}
	return nil, nil, output
}

func anyToArgon(x any, quote bool) string {
	switch x := x.(type) {
	case string:
		if !quote {
			return x
		} else {
			return strconv.Quote(x)
		}
	case float64:
		return fmt.Sprint(x)
	case bool:
		if x {
			return "yes"
		} else {
			return "no"
		}
	case nil:
		return "unknown"
	case []any:
		output := []string{}
		for i := 0; i < len(x); i++ {
			output = append(output, anyToArgon(x[i], true))
		}
		return "[" + strings.Join(output, ", ") + "]"
	case variable:
		return "<function " + x.variable + ">"
	default:
		return fmt.Sprint(x)
	}
}

func runprocess(codeseg any, origin string, vargroups []map[string]variableValue) (any, any) {
	switch codeseg := codeseg.(type) {
	case opperator:
		return runOperator(codeseg, vargroups, origin), "opperator"
	case variable:
		var varigroup map[string]variableValue
		for i := len(vargroups) - 1; i >= 0; i-- {
			if vargroups[i][codeseg.variable].EXISTS != nil {
				varigroup = vargroups[i]
				break
			}
		}
		myvar := varigroup[codeseg.variable]
		if myvar.EXISTS == nil {
			myvar = vargroups[len(vargroups)-1][codeseg.variable]
		}
		if myvar.EXISTS != nil {
			if myvar.TYPE != "func" && myvar.TYPE != "init_function" {
				return myvar.VAL, "value"
			} else {
				return codeseg, "value"
			}
		}
		log.Fatal("undecared variable " + codeseg.variable + " on line " + fmt.Sprint(codeseg.line+1))
	case funcCallType:
		return callFunc(codeseg, vargroups, origin), "value"
	case itemsType:
		vals := []any{}
		for i := 0; i < len(codeseg.vals); i++ {
			resp, _ := runprocess(codeseg.vals[i], origin, vargroups)
			vals = append(vals, resp)
		}
		return vals, "value"
	case whileLoop:
		whileloop := codeseg
		resp, _ := runop(whileloop.condition, origin, vargroups)
		vari := append(vargroups, map[string]variableValue{})
		for boolean(resp) {
			ty, val, _ := run(whileloop.code, origin, vari)
			if ty == "break" {
				break
			} else if ty == "continue" {
				continue
			} else if ty == "return" {
				return val, ty
			}
			resp, _ = runop(whileloop.condition, origin, vargroups)
		}
		return nil, nil
	case ifstatement:
		vari := append(vargroups, map[string]variableValue{})
		iff := codeseg
		for i := 0; i < len(iff.statments); i++ {
			resp, _ := runop(iff.statments[i].condition, origin, vargroups)
			if boolean(resp) {
				ty, val, _ := run(iff.statments[i].code, origin, vari)
				return val, ty
			}
		}
		ty, val, _ := run(iff.FALSE, origin, vari)
		return val, ty
	case importType:
		resp, _ := runop(codeseg.path, origin, vargroups)
		importvars := importMod(resp.(string), filepath.Dir(origin))
		if codeseg.toImport == nil {
			for name, val := range importvars {
				vargroups[len(vargroups)-1][name] = val
			}
		} else {
			var toImport = codeseg.toImport.([]string)
			for i := 0; i < len(toImport); i++ {
				vargroups[len(vargroups)-1][toImport[i]] = importvars[toImport[i]]
			}
		}
		return nil, nil
	case setVariable:
		setVariableVal(codeseg, vargroups, origin)
		return nil, nil
	case setFunction:
		setFunctionVal(codeseg, vargroups, origin)
		return nil, nil
	case returnType:
		var val any = nil
		if codeseg.val != nil {
			vars, worked := runop(codeseg.val, origin, vargroups)
			if worked == nil {
				log.Fatal("return statement must return a value on line " + fmt.Sprint(codeseg.line+1))
			}
			val = vars
		}
		return val, "return"
	case breakType:
		return nil, "break"
	case continueType:
		return nil, "continue"
	}
	return codeseg, "value"
}

func callFunc(call funcCallType, vargroups []map[string]variableValue, origin string) any {
	var variables map[string]variableValue
	for i := len(vargroups) - 1; i >= 0; i-- {
		if vargroups[i][call.name].EXISTS != nil {
			variables = vargroups[i]
			break
		}
	}
	if variables[call.name].EXISTS == nil {
		log.Fatal("undecared function '" + call.name + "' on line " + fmt.Sprint(call.line+1))
	}
	callable := call
	for variables[callable.name].TYPE != "func" && variables[callable.name].TYPE != "init_function" && fmt.Sprint(reflect.TypeOf(variables[callable.name].VAL)) == "main.variable" {
		callable = funcCallType{name: variables[callable.name].VAL.(variable).variable, args: callable.args, line: callable.line}
	}
	for i := len(vargroups) - 1; i >= 0; i-- {
		if vargroups[i][callable.name].EXISTS != nil {
			variables = vargroups[i]
			break
		}
	}
	if variables[callable.name].TYPE != "func" && variables[callable.name].TYPE != "init_function" {
		log.Fatal("'" + call.name + "' is not a function on line " + fmt.Sprint(call.line+1))
	}
	if variables[callable.name].FUNC {

		argvals := []any{}
		for i := 0; i < len(callable.args); i++ {
			resp, _ := runprocess(callable.args[i], origin, vargroups)
			argvals = append(argvals, resp)
		}
		return variables[callable.name].VAL.(func(...any) any)(argvals...)
	} else {
		argvars := make(map[string]variableValue)
		for i := 0; i < len(callable.args); i++ {
			resp, _ := runprocess(callable.args[i], origin, vargroups)
			name := variables[callable.name].VAL.(setFunction).args[i]
			argvars[name] = variableValue{
				VAL:    resp,
				TYPE:   "var",
				EXISTS: true,
				FUNC:   false,
			}
		}
		ty, val, _ := run(variables[callable.name].VAL.(setFunction).code, origin, append(vargroups, modules[variables[callable.name].origin], argvars))
		if ty != "return" && ty != nil {
			log.Fatal(ty, " is not allowed in top level of function on line ", (callable.line + 1))
		}
		return val
	}
}

func setVariableVal(x setVariable, vargroups [](map[string]variableValue), origin string) {
	var variable map[string]variableValue = nil
	if x.TYPE == "preset" {
		for i := len(vargroups) - 1; i >= 0; i-- {
			if vargroups[i][x.variable].EXISTS != nil {
				variable = vargroups[i]
				break
			}
		}
	} else {
		variable = vargroups[len(vargroups)-1]
	}
	if variable[x.variable].EXISTS == nil {
		resp, _ := runop(x.value, origin, vargroups)
		var TYPE = x.TYPE
		if TYPE == "preset" {
			TYPE = "var"
		}
		vargroups[len(vargroups)-1][x.variable] = variableValue{
			TYPE:   TYPE,
			EXISTS: true,
			VAL:    resp,
			origin: origin,
			FUNC:   false,
		}
	} else if variable[x.variable].TYPE == "var" {
		resp, _ := runop(x.value, origin, vargroups)
		var TYPE = x.TYPE
		if TYPE == "preset" {
			TYPE = "var"
		}
		variable[x.variable] = variableValue{
			TYPE:   TYPE,
			EXISTS: variable[x.variable].EXISTS,
			VAL:    resp,
			origin: origin,
			FUNC:   variable[x.variable].FUNC,
		}
	} else {
		log.Fatal("cannot edit " + variable[x.variable].TYPE + " variable on line " + fmt.Sprint(x.line+1))
	}
}

func setFunctionVal(x setFunction, vargroups []map[string]variableValue, origin string) {
	var variable map[string]variableValue
	for i := len(vargroups) - 1; i >= 0; i-- {
		if vargroups[i][x.name].EXISTS != nil {
			variable = vargroups[i]
			break
		}
	}
	if variable[x.name].EXISTS == nil {
		vargroups[len(vargroups)-1][x.name] = variableValue{
			TYPE:   "func",
			EXISTS: true,
			VAL:    x,
			FUNC:   false,
			origin: origin,
		}
	} else if variable[x.name].TYPE == "func" {
		vari := variable[x.name]
		vari.VAL = x
		vari.FUNC = false
	} else {
		log.Fatal("cannot edit " + variable[x.name].TYPE + " variable on line " + fmt.Sprint(x.line+1))
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
	switch x := x.(type) {
	case float64:
		return x
	case float32:
		return float64(x)
	case int:
		return float64(x)
	case bool:
		if x {
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

func runOperator(opperation opperator, vargroups []map[string]variableValue, origin string) any {
	var output any
opperationloop:
	for i := 0; i < len(opperation.vals); i++ {
		switch opperation.t {
		case 0:
			x, _ := runop(opperation.vals[i], origin, vargroups)
			if output == nil {
				output = x
			} else {
				if boolean(output) && boolean(x) {
					output = x
				} else {
					output = false
					break opperationloop
				}
			}
		case 1:
			x, _ := runop(opperation.vals[i], origin, vargroups)
			if boolean(x) {
				output = x
				break opperationloop
			}
		case 2:
			x, _ := runop(opperation.vals[i], origin, vargroups)
			if output == nil {
				output = x
			} else {
				output = xiny(output, x.([]any))
			}
		case 3:
			x, _ := runop(opperation.vals[i], origin, vargroups)
			if output == nil {
				output = x
			} else {
				output = !xiny(output, x.([]interface{}))
			}
		case 4:
			x, _ := runop(opperation.vals[i], origin, vargroups)
			if output == nil {
				output = x
			} else {
				output = (number(output) <= number(x))
			}
		case 5:
			x, _ := runop(opperation.vals[i], origin, vargroups)
			if output == nil {
				output = x
			} else {
				output = (number(output) >= number(x))
			}
		case 6:
			x, _ := runop(opperation.vals[i], origin, vargroups)
			if output == nil {
				output = x
			} else {
				output = (number(output) < number(x))
			}
		case 7:
			x, _ := runop(opperation.vals[i], origin, vargroups)
			if output == nil {
				output = x
			} else {
				output = (number(output) > number(x))
			}
		case 8:
			x, _ := runop(opperation.vals[i], origin, vargroups)
			if output == nil {
				output = x
			} else {
				output = (output != x)
			}
		case 9:
			x, _ := runop(opperation.vals[i], origin, vargroups)
			if output == nil {
				output = x
			} else {
				output = (output == x)
			}
		case 11:
			x, _ := runop(opperation.vals[i], origin, vargroups)
			if output == nil {
				output = x
			} else {
				output = (number(output) - number(x))
			}
		case 10:
			x, _ := runop(opperation.vals[i], origin, vargroups)
			if output == nil {
				output = x
			} else {
				output = dynamicAdd(output, x)
			}
		case 12:
			x, _ := runop(opperation.vals[i], origin, vargroups)
			if output == nil {
				output = x
			} else {
				output = (number(output) * number(x))
			}
		case 14:
			x, _ := runop(opperation.vals[i], origin, vargroups)
			if output == nil {
				output = x
			} else {
				output = math.Floor(number(output) / number(x))
			}
		case 13:
			x, _ := runop(opperation.vals[i], origin, vargroups)
			if output == nil {
				output = x
			} else {
				output = math.Mod(number(output), number(x))
			}
		case 15:
			x, _ := runop(opperation.vals[i], origin, vargroups)
			if output == nil {
				output = x
			} else {
				output = (number(output) / number(x))
			}
		case 16:
			x, _ := runop(opperation.vals[i], origin, vargroups)
			if output == nil {
				output = x
			} else {
				output = math.Pow(number(output), 1/number(x))
			}
		case 17:
			x, _ := runop(opperation.vals[i], origin, vargroups)
			if output == nil {
				output = x
			} else {
				output = math.Pow(number(output), number(x))
			}
		}
	}
	return output
}
