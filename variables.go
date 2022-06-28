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
		TYPE:   "init_function",
		EXISTS: true,
		VAL:    ArgonLog,
		FUNC:   true,
	}
	vars["input"] = variableValue{
		TYPE:   "init_function",
		EXISTS: true,
		VAL:    ArgonInput,
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
	vars["number"] = variableValue{
		TYPE:   "init_function",
		EXISTS: true,
		VAL:    ArgonNumber,
		FUNC:   true,
	}
	vars["whole"] = variableValue{
		TYPE:   "init_function",
		EXISTS: true,
		VAL:    ArgonWhole,
		FUNC:   true,
	}
	vars["string"] = variableValue{
		TYPE:   "init_function",
		EXISTS: true,
		VAL:    ArgonString,
		FUNC:   true,
	}
	vars["license"] = variableValue{
		TYPE:   "init_function",
		EXISTS: true,
		VAL:    ArgonLicense,
		FUNC:   true,
	}
	vars["exec"] = variableValue{
		TYPE:   "init_function",
		EXISTS: true,
		VAL:    exec,
		FUNC:   true,
	}
	vars["eval"] = variableValue{
		TYPE:   "init_function",
		EXISTS: true,
		VAL:    eval,
		FUNC:   true,
	}
}
