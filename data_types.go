package main

type process struct {
	t    int
	x    interface{}
	y    interface{}
	line int
}

type variable struct {
	variable string
	line     int
}

type whileLoop struct {
	condition interface{}
	code      []code
}

type setVariable struct {
	variable string
	value    interface{}
	line     int
}
