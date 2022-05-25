package main

type variableValue struct {
	TYPE string
	VAL  interface{}
	FUNC bool
}

var builtins = make(map[string]variableValue)

func init() {
	builtins["cool"] = variableValue{
		TYPE: "string",
		VAL:  "hello world",
		FUNC: false,
	}
}
