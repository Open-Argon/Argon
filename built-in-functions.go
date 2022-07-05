package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

func ArgonLog(x ...any) (any, any) {
	output := []any{}
	for i := 0; i < len(x); i++ {
		output = append(output, anyToArgon(x[i], false))
	}
	fmt.Println(output...)
	return nil, nil
}

func ArgonNumber(x ...any) (any, any) {
	return (number(x[0]))
}

func ArgonString(x ...any) (any, any) {
	return (anyToArgon(x[0], true)), nil
}

func ArgonWhole(x ...any) (any, any) {
	num, err := number(x[0])
	if err != nil {
		return nil, err
	}
	return (math.Floor(num)), nil
}

func ArgonInput(x ...any) (any, any) {
	fmt.Print(x[0])
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	line := scanner.Text()
	if err := scanner.Err(); err != nil {
		fmt.Println("Error encountered:", err)
	}
	return (line), nil
}

func ArgonLicense(...any) (any, any) {
	fmt.Println(`MIT License

Copyright (c) 2022 William Bell

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.`)
	return nil, nil
}

func exec(x ...any) (any, any) {
	origin := ""
	if len(x) > 1 {
		origin = x[1].(string)
	}
	runStr(x[0].(string), origin, map[string]variableValue{})
	return nil, nil
}

func eval(x ...any) (any, any) {
	resp, _, err := translateprocess(code{
		code: x[0].(string),
		line: 0,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ArgonAppend(x ...any) (any, any) {
	return append(x[0].([]any), x[1]), nil
}

func ArgonExtend(x ...any) (any, any) {
	return append(x[0].([]any), x[1].([]any)...), nil
}

func ArgonLen(x ...any) (any, any) {
	switch value := x[0].(type) {
	case []any:
		return len(value), nil
	case string:
		return len(value), nil
	case map[any]any:
		return len(value), nil
	default:
		return 0, nil
	}
}

func ArgonJoin(x ...any) (any, any) {
	return strings.Join(x[0].([]string), x[1].(string)), nil
}

func ArgonTime(x ...any) (any, any) {
	return time.Now().UnixMilli(), nil
}

func ArgonSleep(x ...any) (any, any) {
	num, err := number(x[0])
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Duration(num) * time.Millisecond)
	return nil, nil
}

func ArgonSetSeed(x ...any) (any, any) {
	rand.Seed(x[0].(int64))
	return nil, nil
}

func ArgonRandom(x ...any) (any, any) {
	return rand.Float64(), nil
}
