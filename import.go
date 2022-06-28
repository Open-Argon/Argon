package main

import (
	"log"
	"os"
	"path/filepath"
)

var imported = []any{}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func importMod(path string) {
	extention := filepath.Ext(path)
	if extention == "" {
		path += ".ar"
	}
	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	exPath := ex
	p := filepath.Join(exPath, path)
	if !FileExists(p) {
		ex, err = os.Executable()
		if err != nil {
			panic(err)
		}
		exPath = filepath.Dir(ex)
		p = filepath.Join(exPath, "modules", path)
		if !FileExists(p) {
			log.Fatal("Module not found: " + path + " (" + p + ")")
		}
	}
	if xiny(p, imported) {
		return
	}
	imported = append(imported, p)
	data, err := os.ReadFile(p)
	if err != nil {
		log.Fatal(err)
	}
	runStr(string(data))
}

func runStr(str string) [][]any {
	translated := translate(str)
	ty, _, resp := run(translated, make(map[string]variableValue))
	if ty != nil {
		log.Fatal(ty, " at top level")
	}
	return resp
}
