package main

import (
	"log"
	"os"
	"path/filepath"
)

var imported = []any{}

func importMod(path string) {
	if xiny(path, imported) {
		return
	}
	imported = append(imported, path)
	extention := filepath.Ext(path)
	if extention == "" {
		path += ".ar"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	script := translate(string(data))
	run(script)
}
