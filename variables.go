package main

import "math"

var vars = make(map[string]variableValue)

func init() {
	vars["log"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    builtinFunc{"log", ArgonLog},
		FUNC:   true,
	}
	vars["input"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    builtinFunc{"input", ArgonInput},
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
	vars["infinity"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    math.Inf(1),
		FUNC:   false,
	}
	vars["NaN"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    math.NaN(),
		FUNC:   false,
	}
	vars["unknown"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    nil,
		FUNC:   false,
	}
	vars["number"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    builtinFunc{"number", ArgonNumber},
		FUNC:   true,
	}
	vars["whole"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    builtinFunc{"whole", ArgonWhole},
		FUNC:   true,
	}
	vars["string"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    builtinFunc{"string", ArgonString},
		FUNC:   true,
	}
	vars["license"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    builtinFunc{"license", ArgonLicense},
		FUNC:   true,
	}
	vars["exec"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    builtinFunc{"exec", exec},
		FUNC:   true,
	}
	vars["eval"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    builtinFunc{"eval", eval},
		FUNC:   true,
	}

	vars["append"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    builtinFunc{"append", ArgonAppend},
		FUNC:   true,
	}
	vars["extend"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    ArgonExtend,
		FUNC:   true,
	}
	vars["len"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    builtinFunc{"len", ArgonLen},
		FUNC:   true,
	}
	vars["join"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    builtinFunc{"join", ArgonJoin},
		FUNC:   true,
	}
	vars["random"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    builtinFunc{"random", ArgonRandom},
		FUNC:   true,
	}
	vars["randomSetSeed"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    builtinFunc{"randomSetSeed", ArgonSetSeed},
		FUNC:   true,
	}
	vars["time"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    builtinFunc{"time", ArgonTime},
		FUNC:   true,
	}
	vars["sleep"] = variableValue{
		TYPE:   "init",
		EXISTS: true,
		VAL:    builtinFunc{"sleep", ArgonSleep},
		FUNC:   true,
	}
}
