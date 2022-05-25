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

type setVariable struct {
	variable string
	value    interface{}
	line     int
}
