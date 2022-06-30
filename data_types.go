package main

type opperator struct {
	t    int
	vals []interface{}
	line int
}

type variable struct {
	variable string
	line     int
}

type funcCallType struct {
	name string
	args []interface{}
	line int
}

type whileLoop struct {
	condition interface{}
	code      []interface{}
}

type variableValue struct {
	TYPE   string
	VAL    interface{}
	EXISTS interface{}
	origin string
	FUNC   bool
}

type importType struct {
	path     interface{}
	toImport any
	line     int
}

type ifstatement struct {
	condition interface{}
	TRUE      []interface{}
	FALSE     []interface{}
}

type setVariable struct {
	TYPE     string
	variable string
	value    interface{}
	line     int
}

type setFunction struct {
	name string
	args []string
	code []interface{}
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
