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

func run(lines []interface{}) {
	fmt.Println(runprocess(opperator{
		t: 11,
		x: opperator{
			t: 11, x: opperator{
				t: 11,
				x: 1e100,
				y: 5,
			}, y: " "},
		y: variable{
			variable: "HW",
			line:     10,
		},
	}))
}

func runprocess(codeseg interface{}) interface{} {
	switch codeseg.(type) {
	case opperator:
		return runOperator(codeseg.(opperator))
	case variable:
		myvar := vars[codeseg.(variable).variable]
		if myvar.EXISTS != nil {
			return myvar.VAL
		}
	}
	return codeseg
}

func dynamicAdd(x interface{}, y interface{}) interface{} {
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
	xnum, err := strconv.ParseFloat(fmt.Sprint(x), 64)
	if err != nil {
		log.Fatal(err)
	}
	ynum, err := strconv.ParseFloat(fmt.Sprint(y), 64)
	if err != nil {
		log.Fatal(err)
	}
	return xnum + ynum
}

func xiny(x interface{}, y []interface{}) bool {
	for i := 0; i < len(y); i++ {
		if x == y[i] {
			return true
		}
	}
	return false
}

func runOperator(opperation opperator) interface{} {
	x := runop(opperation.x)
	y := runop(opperation.y)
	var output interface{}
	switch opperation.t {
	case 0:
		output = (x != false && x != nil) && (y != false && y != nil)
	case 1:
		output = x.(bool) || y.(bool)
	case 2:
		output = xiny(x, y.([]interface{}))
	case 3:
		output = !xiny(x, y.([]interface{}))
	case 4:
		output = (x.(float64) <= y.(float64))
	case 5:
		output = (x.(float64) >= y.(float64))
	case 6:
		output = (x.(float64) < y.(float64))
	case 7:
		output = (x.(float64) > y.(float64))
	case 8:
		output = (x != y)
	case 9:
		output = (x == y)
	case 10:
		output = (x.(float64) - y.(float64))
	case 11:
		output = dynamicAdd(x, y)
	case 12:
		output = math.Pow(x.(float64), 1/y.(float64))
	case 13:
		output = math.Pow(x.(float64), y.(float64))
	case 14:
		output = x.(float64) * y.(float64)
	case 15:
		output = math.Floor(x.(float64) / y.(float64))
	case 16:
		output = math.Floor(x.(float64) / y.(float64))
	}
	return output
}
