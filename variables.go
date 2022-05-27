package main

import (
	"fmt"
)

type variableValue struct {
	TYPE   string
	VAL    interface{}
	EXISTS interface{}
	FUNC   bool
}

var vars = make(map[string]variableValue)

func ArgonLog(x ...any) any {
	fmt.Println(x...)
	return nil
}

func init() {
	vars["log"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    ArgonLog,
		FUNC:   false,
	}
}
