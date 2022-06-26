package main

import (
	"log"
	"os"
	"path/filepath"
)

var imported = []any{}

func importMod(path string) {
	extention := filepath.Ext(path)
	if extention == "" {
		path += ".ar"
	}
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	path = filepath.Join(exPath, path)
	if xiny(path, imported) {
		return
	}
	imported = append(imported, path)
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	runStr(string(data))
}

func runStr(str string) [][]any {
	translated := translate(str)
	return run(translated)
}
