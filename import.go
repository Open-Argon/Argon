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

func importMod(realpath string, origin string) {
	extention := filepath.Ext(realpath)
	path := realpath
	if extention == "" {
		path += ".ar"
	}
	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	executable = filepath.Dir(executable)
	pathsToTest := []string{filepath.Join(origin, realpath, "init.ar"), filepath.Join(origin, path), filepath.Join(origin, "modules", path), filepath.Join(origin, "modules", realpath, "init.ar"), filepath.Join(ex, path), filepath.Join(ex, "modules", realpath, "init.ar"), filepath.Join(ex, "modules", path), filepath.Join(executable, "modules", realpath, "init.ar"), filepath.Join(executable, "modules", path)}
	var p string
	for _, p = range pathsToTest {
		if FileExists(p) {
			break
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
	runStr(string(data), filepath.Dir(p))
}

func runStr(str string, origin string) [][]any {
	translated := translate(str)
	ty, _, resp := run(translated, origin, []map[string]variableValue{vars})
	if ty != nil {
		log.Fatal(ty, " at top level")
	}
	return resp
}
