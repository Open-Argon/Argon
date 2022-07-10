package main

import (
	"fmt"
	"math"
	"path/filepath"
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
			if ty == "return" || ty == "break" || ty == "continue" || ty == "error" {
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
		if math.IsNaN(x) {
			return "NaN"
		} else if math.IsInf(x, 1) {
			return "infinity"
		} else if math.IsInf(x, -1) {
			return "-infinity"
		} else {
			return strconv.FormatFloat(x, 'f', -1, 64)
		}
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
	case builtinFunc:
		return fmt.Sprint("<function ", x.name, ">")
	case callableFunction:
		if x.name != nil {
			return fmt.Sprint("<function ", x.name, ">")
		} else {
			return "<anonymous function>"
		}
	default:
		return fmt.Sprint(x)
	}
}

func runprocess(codeseg any, origin string, vargroups []map[string]variableValue) (any, any) {
	switch codeseg := codeseg.(type) {
	case opperator:
		resp, ty := runOperator(codeseg, vargroups, origin)
		return resp, ty
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
			return myvar.VAL, "value"
		}
		return ("undecared variable " + codeseg.variable + ": " + origin + ":" + fmt.Sprint(codeseg.line+1)), "error"
	case funcCallType:
		resp, ty := callFunc(codeseg, vargroups, origin)
		return resp, ty
	case errorType:
		resp, _ := runprocess(codeseg.val, origin, vargroups)
		return resp, "error"
	case itemsType:
		vals := []any{}
		for i := 0; i < len(codeseg.vals); i++ {
			resp, err := runprocess(codeseg.vals[i], origin, vargroups)
			if err != nil {
				return "invalid value: " + origin + ":" + fmt.Sprint(codeseg.line+1), "error"
			}
			vals = append(vals, resp)
		}
		return vals, "value"
	case tryType:

		ty, resp, _ := run(codeseg.code, origin, append(vargroups, map[string]variableValue{}))
		if ty == "error" {
			ty, resp, _ := run(codeseg.catch, origin, append(vargroups, map[string]variableValue{"err": {
				TYPE:   "var",
				EXISTS: true,
				VAL:    resp,
				FUNC:   false,
			}}))
			return resp, ty
		}
		return resp, ty
	case whileLoop:
		whileloop := codeseg
		resp, ty := runop(whileloop.condition, origin, vargroups)
		if ty == "error" {
			return resp, ty
		}
		vari := append(vargroups, map[string]variableValue{})
		for boolean(resp) {
			ty, val, _ := run(whileloop.code, origin, vari)
			if ty == "break" {
				break
			} else if ty == "continue" {
				continue
			} else if ty == "return" || ty == "error" {
				return val, ty
			}
			resp, ty = runop(whileloop.condition, origin, vargroups)
			if ty == "error" {
				return resp, ty
			}
		}
		return nil, nil
	case ifstatement:
		vari := append(vargroups, map[string]variableValue{})
		iff := codeseg
		for i := 0; i < len(iff.statments); i++ {
			resp, ty := runop(iff.statments[i].condition, origin, vargroups)
			if ty == "error" {
				return resp, ty
			}
			if boolean(resp) {
				ty, val, _ := run(iff.statments[i].code, origin, vari)
				return val, ty
			}
		}
		ty, val, _ := run(iff.FALSE, origin, vari)
		return val, ty
	case importType:
		resp, ty := runop(codeseg.path, origin, vargroups)
		if ty == "error" {
			return resp, ty
		}
		importvars, err := importMod(fmt.Sprint(resp), filepath.Dir(origin), false)
		if err != nil {
			return err, "error"
		}
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
		value := setVariableVal(codeseg, vargroups, origin)
		if value == nil {
			return nil, nil
		} else {
			return value, "error"
		}
	case setFunction:
		value := setFunctionVal(codeseg, vargroups, origin)
		if value == nil {
			return nil, nil
		} else {
			return value, "error"
		}
	case indexType:
		to, ty := runop(codeseg.to, origin, vargroups)
		if ty == "error" {
			return to, ty
		}

		index, ty := runop(codeseg.index, origin, vargroups)
		if ty == "error" {
			return to, ty
		}
		switch value := to.(type) {
		case string:
			index, err := number(index)
			if err != nil {
				return err, "error"
			}
			i := int(index)
			if i >= len(value) {
				return "index out of range on line " + origin + ":" + fmt.Sprint(codeseg.line+1), "error"
			} else if i < 0 {
				i = len(value) + i
			}
			return string(value[i]), "value"
		case []any:
			index, err := number(index)
			if err != nil {
				return err, "error"
			}
			i := int(index)
			if i >= len(value) {
				return "index out of range on line " + origin + ":" + fmt.Sprint(codeseg.line+1), "error"
			} else if i < 0 {
				i = len(value) + i
			}
			return value[i], "value"
		case map[string]any:
			key := fmt.Sprint(index)
			if value[key] == nil {
				return "'" + key + "' not in value on line " + origin + ":" + fmt.Sprint(codeseg.line+1), "error"
			}
			return value[key], "value"
		default:
			return "value not indexable " + origin + ":" + fmt.Sprint(codeseg.line+1), "error"
		}
	case returnType:
		var val any = nil
		var ty any = "return"
		if codeseg.val != nil {
			vars, worked := runop(codeseg.val, origin, vargroups)
			if worked == nil {
				return "return statement must return a value: " + origin + ":" + fmt.Sprint(codeseg.line+1), "error"
			} else if worked == "error" {
				return vars, "error"
			}
			val = vars
		}
		return val, ty
	case deleteType:
		var varigroup map[string]variableValue
		for i := len(vargroups) - 1; i >= 0; i-- {
			if vargroups[i][codeseg.variable.variable].EXISTS != nil {
				varigroup = vargroups[i]
				break
			}
		}
		if varigroup[codeseg.variable.variable].EXISTS == nil {
			varigroup = vargroups[len(vargroups)-1]
		}
		if varigroup[codeseg.variable.variable].EXISTS != nil {
			delete(varigroup, codeseg.variable.variable)
			return nil, nil
		}
		return "undeclared variable " + codeseg.variable.variable + ": " + origin + ":" + fmt.Sprint(codeseg.line+1), "error"
	case breakType:
		return nil, "break"
	case continueType:
		return nil, "continue"
	}
	return codeseg, "value"
}

func callFunc(call funcCallType, vargroups []map[string]variableValue, origin string) (any, any) {
	resp, ty := runprocess(call.FUNC, origin, vargroups)
	if ty == "error" {
		return resp, ty
	}

	switch variable := resp.(type) {
	case callableFunction:
		if len(call.args) != len(variable.args) {
			return "invalid number of arguments: " + origin + ":" + fmt.Sprint(call.line+1), "error"
		}

		argvars := make(map[string]variableValue)
		for i := 0; i < len(variable.args); i++ {
			resp, ty := runprocess(call.args[i], origin, vargroups)
			if ty == "error" {
				return resp, ty
			}
			name := variable.args[i]
			argvars[name] = variableValue{
				VAL:    resp,
				TYPE:   "var",
				EXISTS: true,
				FUNC:   false,
			}
		}
		ty, val, _ := run(variable.code, origin, append(variable.vargroups, argvars))
		if ty != "return" && ty != "error" && ty != nil {
			return fmt.Sprint(ty) + " is not allowed in function: " + origin + ":" + fmt.Sprint(variable.line+1), "error"
		} else if ty == "error" {
			return val, ty
		}
		return val, "value"
	case builtinFunc:
		argvals := []any{}
		for i := 0; i < len(call.args); i++ {
			resp, ty := runprocess(call.args[i], origin, vargroups)
			if ty == "error" {
				return resp, ty
			}
			argvals = append(argvals, resp)
		}
		val, err := variable.FUNC(argvals...)
		if err != nil {
			return err, "error"
		}
		return val, "value"
	default:
		return ("value is not a function: " + origin + ":" + fmt.Sprint(call.line+1)), "error"
	}
}

func isCallableFunc(x any) bool {
	switch x.(type) {
	case callableFunction:
		return true
	default:
		return false
	}
}

func setVariableVal(x setVariable, vargroups [](map[string]variableValue), origin string) any {
	var variable map[string]variableValue = nil
	if x.TYPE == "preset" {
		for i := len(vargroups) - 1; i >= 0; i-- {
			if vargroups[i][x.variable.variable].EXISTS != nil {
				variable = vargroups[i]
				break
			}
		}
	} else {
		variable = vargroups[len(vargroups)-1]
	}
	if variable[x.variable.variable].EXISTS == nil {
		var val any
		var TYPE = x.TYPE
		if isCallableFunc(x.value) {
			val = x.value
		} else {
			resp, ty := runop(x.value, origin, vargroups)
			if ty == "error" {
				return resp
			}
			val = resp
			if TYPE == "preset" {
				TYPE = "var"
			}
		}
		vargroups[len(vargroups)-1][x.variable.variable] = variableValue{
			TYPE:   TYPE,
			EXISTS: true,
			VAL:    val,
			origin: origin,
			FUNC:   false,
		}
	} else if variable[x.variable.variable].TYPE == "var" {
		resp, ty := runop(x.value, origin, vargroups)
		if ty == "error" {
			return resp
		}
		var TYPE = x.TYPE
		if TYPE == "preset" {
			TYPE = "var"
		}
		variable[x.variable.variable] = variableValue{
			TYPE:   TYPE,
			EXISTS: variable[x.variable.variable].EXISTS,
			VAL:    resp,
			origin: origin,
			FUNC:   variable[x.variable.variable].FUNC,
		}
	} else {
		return ("cannot edit " + variable[x.variable.variable].TYPE + " variable: " + origin + ":" + fmt.Sprint(x.line+1))
	}
	return nil
}

func setFunctionVal(x setFunction, vargroups []map[string]variableValue, origin string) any {
	var variable map[string]variableValue
	for i := len(vargroups) - 1; i >= 0; i-- {
		if vargroups[i][x.name].EXISTS != nil {
			variable = vargroups[i]
			break
		}
	}
	x.FUNC.vargroups = vargroups
	if variable[x.name].EXISTS == nil {
		vargroups[len(vargroups)-1][x.name] = variableValue{
			TYPE:   "var",
			EXISTS: true,
			VAL:    x.FUNC,
			FUNC:   false,
			origin: origin,
		}
	} else if variable[x.name].TYPE == "var" {
		vari := variable[x.name]
		vari.VAL = x
		vari.FUNC = false
	} else {
		return ("cannot edit " + variable[x.name].TYPE + " variable: " + origin + ":" + fmt.Sprint(x.line+1))
	}
	return nil
}

func dynamicAdd(x any, y any) (any, any) {
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
		return fmt.Sprint(x) + fmt.Sprint(y), nil
	}
	xnum, err := number(x)
	if err != nil {
		return err, "error"
	}
	ynum, err := number(y)
	if err != nil {
		return err, "error"
	}
	return xnum + ynum, nil
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

func number(x any) (float64, any) {
	switch x := x.(type) {
	case float64:
		return x, nil
	case float32:
		return float64(x), nil
	case int:
		return float64(x), nil
	case bool:
		if x {
			return 1, nil
		} else {
			return 0, nil
		}
	default:
		num, err := strconv.ParseFloat(fmt.Sprint(x), 64)
		if err != nil {
			return 0, (err)
		}
		return num, nil
	}
}

func runOperator(opperation opperator, vargroups []map[string]variableValue, origin string) (any, any) {
	var output any
opperationloop:
	for i := 0; i < len(opperation.vals); i++ {
		switch opperation.t {
		case 0:
			x, ty := runop(opperation.vals[i], origin, vargroups)
			if ty == "error" {
				return x, ty
			}
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
			x, ty := runop(opperation.vals[i], origin, vargroups)
			if ty == "error" {
				return x, ty
			}
			if boolean(x) {
				output = x
				break opperationloop
			}
		case 2:
			x, ty := runop(opperation.vals[i], origin, vargroups)
			if ty == "error" {
				return x, ty
			}
			if output == nil {
				output = x
			} else {
				output = xiny(output, x.([]any))
			}
		case 3:
			x, ty := runop(opperation.vals[i], origin, vargroups)
			if ty == "error" {
				return x, ty
			}
			if output == nil {
				output = x
			} else {
				output = !xiny(output, x.([]interface{}))
			}
		case 4:
			x, ty := runop(opperation.vals[i], origin, vargroups)
			if ty == "error" {
				return x, ty
			}
			if output == nil {
				output = x
			} else {
				out, err := number(output)
				if err != nil {
					return 0, (err)
				}

				in, err := number(x)
				if err != nil {
					return 0, (err)
				}
				output = (out <= in)
			}
		case 5:
			x, ty := runop(opperation.vals[i], origin, vargroups)
			if ty == "error" {
				return x, ty
			}
			if output == nil {
				output = x
			} else {
				out, err := number(output)
				if err != nil {
					return 0, (err)
				}

				in, err := number(x)
				if err != nil {
					return 0, (err)
				}
				output = (out >= in)
			}
		case 6:
			x, ty := runop(opperation.vals[i], origin, vargroups)
			if ty == "error" {
				return x, ty
			}
			if output == nil {
				output = x
			} else {
				out, err := number(output)
				if err != nil {
					return 0, (err)
				}

				in, err := number(x)
				if err != nil {
					return 0, (err)
				}
				output = (out < in)
			}
		case 7:
			x, ty := runop(opperation.vals[i], origin, vargroups)
			if ty == "error" {
				return x, ty
			}
			if output == nil {
				output = x
			} else {
				out, err := number(output)
				if err != nil {
					return 0, (err)
				}

				in, err := number(x)
				if err != nil {
					return 0, (err)
				}
				output = (out > in)
			}
		case 8:
			x, ty := runop(opperation.vals[i], origin, vargroups)
			if ty == "error" {
				return x, ty
			}
			if output == nil {
				output = x
			} else {
				output = (output != x)
			}
		case 9:
			x, ty := runop(opperation.vals[i], origin, vargroups)
			if ty == "error" {
				return x, ty
			}
			if output == nil {
				output = x
			} else {
				output = (output == x)
			}
		case 11:
			x, ty := runop(opperation.vals[i], origin, vargroups)
			if ty == "error" {
				return x, ty
			}
			if output == nil {
				output = x
			} else {
				out, err := number(output)
				if err != nil {
					return 0, (err)
				}

				in, err := number(x)
				if err != nil {
					return 0, (err)
				}
				output = (out - in)
			}
		case 10:
			x, ty := runop(opperation.vals[i], origin, vargroups)
			if ty == "error" {
				return x, ty
			}
			if output == nil {
				output = x
			} else {
				val, err := dynamicAdd(output, x)
				if err != nil {
					return 0, (err)
				}
				output = val
			}
		case 12:
			x, ty := runop(opperation.vals[i], origin, vargroups)
			if ty == "error" {
				return x, ty
			}
			if output == nil {
				output = x
			} else {
				out, err := number(output)
				if err != nil {
					return 0, (err)
				}

				in, err := number(x)
				if err != nil {
					return 0, (err)
				}
				output = (out * in)
			}
		case 14:
			x, ty := runop(opperation.vals[i], origin, vargroups)
			if ty == "error" {
				return x, ty
			}
			if output == nil {
				output = x
			} else {
				out, err := number(output)
				if err != nil {
					return 0, (err)
				}

				in, err := number(x)
				if err != nil {
					return 0, (err)
				}
				output = math.Floor(out / in)
			}
		case 13:
			x, ty := runop(opperation.vals[i], origin, vargroups)
			if ty == "error" {
				return x, ty
			}
			if output == nil {
				output = x
			} else {
				out, err := number(output)
				if err != nil {
					return 0, (err)
				}

				in, err := number(x)
				if err != nil {
					return 0, (err)
				}
				output = math.Mod(out, in)
			}
		case 15:
			x, ty := runop(opperation.vals[i], origin, vargroups)
			if ty == "error" {
				return x, ty
			}
			if output == nil {
				output = x
			} else {
				out, err := number(output)
				if err != nil {
					return 0, (err)
				}

				in, err := number(x)
				if err != nil {
					return 0, (err)
				}
				output = (out / in)
			}
		case 16:
			x, ty := runop(opperation.vals[i], origin, vargroups)
			if ty == "error" {
				return x, ty
			}
			if output == nil {
				output = x
			} else {
				out, err := number(output)
				if err != nil {
					return 0, (err)
				}

				in, err := number(x)
				if err != nil {
					return 0, (err)
				}
				output = math.Pow(out, 1/in)
			}
		case 17:
			x, ty := runop(opperation.vals[i], origin, vargroups)
			if ty == "error" {
				return x, ty
			}
			if output == nil {
				output = x
			} else {
				out, err := number(output)
				if err != nil {
					return 0, (err)
				}

				in, err := number(x)
				if err != nil {
					return 0, (err)
				}
				output = math.Pow(out, in)
			}
		}
	}
	return output, "opperator"
}
