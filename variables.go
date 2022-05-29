package main

type variableValue struct {
	TYPE   string
	VAL    interface{}
	EXISTS interface{}
	FUNC   bool
}

var vars = make(map[string]variableValue)

func init() {
	vars["log"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    ArgonLog,
		FUNC:   true,
	}
	vars["number"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    ArgonNumber,
		FUNC:   true,
	}
	vars["yes"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    true,
		FUNC:   false,
	}
	vars["no"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    false,
		FUNC:   false,
	}
	vars["unknown"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    nil,
		FUNC:   false,
	}
}
