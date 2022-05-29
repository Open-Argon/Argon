package main

import (
	"fmt"
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
	fmt.Println(path)
	if xiny(path, imported) {
		return
	}
	imported = append(imported, path)
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	script := translate(string(data))
	run(script)
}
