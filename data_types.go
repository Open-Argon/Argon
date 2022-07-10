package main

type opperator struct {
	t    int
	vals []any
	line int
}

type variable struct {
	variable string
	line     int
}

type funcCallType struct {
	FUNC any
	args []any
	line int
}

type whileLoop struct {
	condition any
	code      []any
}

type variableValue struct {
	TYPE   string
	VAL    any
	EXISTS any
	origin string
	FUNC   bool
}

type importType struct {
	path     any
	toImport any
	line     int
}

type iftype struct {
	condition any
	code      []any
}

type ifstatement struct {
	statments []iftype
	FALSE     []any
}

type setVariable struct {
	TYPE     string
	variable variable
	value    any
	line     int
}

type setFunction struct {
	name string
	FUNC callableFunction
	line int
}

type returnType struct {
	val  any
	line int
}

type breakType struct {
	line int
}

type continueType struct {
	line int
}

type itemsType struct {
	vals []any
	line int
}

type tryType struct {
	code  []any
	catch []any
	line  int
}

type errorType struct {
	val any
}

type deleteType struct {
	variable variable
	line     int
}

type builtinFunc struct {
	name string
	FUNC func(...any) (any, any)
}

type callableFunction struct {
	name      any
	args      []string
	code      []any
	vargroups []map[string]variableValue
	line      int
}

type indexType struct {
	to    any
	index any
	line  int
}
