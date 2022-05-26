package main

type variableValue struct {
	TYPE   string
	VAL    interface{}
	EXISTS interface{}
	FUNC   bool
}

var vars = make(map[string]variableValue)

func init() {
	vars["HW"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    "hello world",
		FUNC:   false,
	}
}
