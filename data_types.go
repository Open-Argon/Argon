package main

type opperator struct {
	t    int
	x    interface{}
	y    interface{}
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
type importType struct {
	path interface{}
	line int
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
