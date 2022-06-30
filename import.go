package main

import (
	"log"
	"os"
	"path/filepath"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

var modules = map[string](map[string]variableValue){}

func importMod(realpath string, origin string) map[string]variableValue {
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
	if modules[p] == nil {
		modules[p] = map[string]variableValue{}
		data, err := os.ReadFile(p)
		if err != nil {
			log.Fatal(err)
		}
		runStr(string(data), p, modules[p])
	}
	return modules[p]
}

func runStr(str string, origin string, variables map[string]variableValue) [][]any {
	translated := translate(str)
	ty, _, resp := run(translated, origin, []map[string]variableValue{vars, variables})
	if ty != nil {
		log.Fatal(ty, " at top level")
	}
	return resp
}
