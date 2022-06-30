package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
)

func ArgonLog(x ...any) any {
	output := []any{}
	for i := 0; i < len(x); i++ {
		output = append(output, anyToArgon(x[i], false))
	}
	fmt.Println(output...)
	return nil
}

func ArgonNumber(x ...any) any {
	return (number(x[0]))
}

func ArgonString(x ...any) any {
	return (anyToArgon(x[0], true))
}

func ArgonWhole(x ...any) any {
	return (math.Floor(number(x[0])))
}

func ArgonInput(x ...any) any {
	fmt.Print(x[0])
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	line := scanner.Text()
	if err := scanner.Err(); err != nil {
		fmt.Println("Error encountered:", err)
	}
	return (line)
}

func ArgonLicense(...any) any {
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
	return nil
}

func exec(x ...any) any {
	origin := ""
	if len(x) > 1 {
		origin = x[1].(string)
	}
	runStr(x[0].(string), origin)
	return nil
}

func eval(x ...any) any {
	resp, _ := translateprocess(code{
		code: x[0].(string),
		line: 0,
	})
	return resp
}

func ArgonAppend(x ...any) any {
	return append(x[0].([]any), x[1])
}

func ArgonExtend(x ...any) any {
	return append(x[0].([]any), x[1].([]any)...)
}

func ArgonLen(x ...any) any {
	switch value := x[0].(type) {
	case []any:
		return len(value)
	case string:
		return len(value)
	case map[any]any:
		return len(value)
	default:
		return 0
	}
}

func ArgonJoin(x ...any) any {
	return strings.Join(x[0].([]string), x[1].(string))
}

func ArgonSetSeed(x ...any) any {
	rand.Seed(x[0].(int64))
	return nil
}

func ArgonRandom(x ...any) any {
	return rand.Float64()
}
