package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
)

var runop func(codeseg interface{}) interface{}

func init() {
	runop = runprocess
}

func run(lines []any) {
	for i := 0; i < len(lines); i++ {
		if lines[i] != nil {
			fmt.Println(runprocess(lines[i]))
		}
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
		log.Fatal("undecared variable on line " + fmt.Sprint(codeseg.(variable).line))
	case setVariable:

	}
	return codeseg
}

func setVariableVal(x setVariable) {
	variable := vars[x.variable]
	if variable.EXISTS != nil || variable.TYPE == "var" {
		vars[x.variable] = variableValue{
			TYPE:   x.variable,
			EXISTS: true,
			VAL:    x.value,
			FUNC:   false,
		}
	} else if variable.EXISTS != nil {
		log.Fatal("undecared variable on line " + fmt.Sprint(x.line))
	} else {
		log.Fatal("cannot edit " + variable.TYPE + " variable on line " + fmt.Sprint(x.line))
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

func number(x any) float64 {
	num, err := strconv.ParseFloat(fmt.Sprint(x), 64)
	if err != nil {
		log.Fatal(err)
	}
	return num
}

func runOperator(opperation opperator) any {
	var output any
	switch opperation.t {
	case 0:
		x := runop(opperation.x)
		y := runop(opperation.y)
		if (x != false && x != nil) && (y != false && y != nil) {
			output = y
		} else {
			output = false
		}
	case 1:
		x := runop(opperation.x)
		if x != false && x != nil {
			output = x
		} else {
			y := runop(opperation.y)
			if y != false && y != nil {
				output = y
			} else {
				output = false
			}
		}
	case 2:
		x := runop(opperation.x)
		y := runop(opperation.y)
		output = xiny(x, y.([]interface{}))
	case 3:
		x := runop(opperation.x)
		y := runop(opperation.y)
		output = !xiny(x, y.([]interface{}))
	case 4:
		x := runop(opperation.x)
		y := runop(opperation.y)
		output = (number(x) <= number(y))
	case 5:
		x := runop(opperation.x)
		y := runop(opperation.y)
		output = (number(x) >= number(y))
	case 6:
		x := runop(opperation.x)
		y := runop(opperation.y)
		output = (number(x) < number(y))
	case 7:
		x := runop(opperation.x)
		y := runop(opperation.y)
		output = (number(x) > number(y))
	case 8:
		x := runop(opperation.x)
		y := runop(opperation.y)
		output = (x != y)
	case 9:
		x := runop(opperation.x)
		y := runop(opperation.y)
		output = (x == y)
	case 10:
		x := runop(opperation.x)
		y := runop(opperation.y)
		output = (number(x) - number(y))
	case 11:
		x := runop(opperation.x)
		y := runop(opperation.y)
		output = dynamicAdd(x, y)
	case 12:
		x := runop(opperation.x)
		y := runop(opperation.y)
		output = math.Pow(number(x), 1/number(y))
	case 13:
		x := runop(opperation.x)
		y := runop(opperation.y)
		output = math.Pow(number(x), number(y))
	case 14:
		x := runop(opperation.x)
		y := runop(opperation.y)
		output = number(x) * number(y)
	case 15:
		x := runop(opperation.x)
		y := runop(opperation.y)
		output = math.Floor(number(x) / number(y))
	case 16:
		x := runop(opperation.x)
		y := runop(opperation.y)
		output = math.Floor(number(x) / number(y))
	}
	return output
}
