package main

import "fmt"

func ArgonLog(x ...any) any {
	output := []any{}
	for i := 0; i < len(x); i++ {
		output = append(output, anyToArgon(x[i]))
	}
	fmt.Println(output...)
	return nil
}

func ArgonNumber(x ...any) any {
	return (number(x[0]))
}
