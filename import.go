package main

import (
	"fmt"
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
var importing = map[string]any{}

func importMod(realpath string, origin string, main bool) (map[string]variableValue, any) {
	extention := filepath.Ext(realpath)
	path := realpath
	if extention == "" {
		path += ".ar"
	}
	ex, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	executable, err := os.Executable()
	if err != nil {
		return nil, err
	}
	executable = filepath.Dir(executable)
	pathsToTest := []string{filepath.Join(origin, realpath, "init.ar"), filepath.Join(origin, path), filepath.Join(origin, "modules", path), filepath.Join(origin, "modules", realpath, "init.ar"), filepath.Join(ex, path), filepath.Join(ex, "modules", realpath, "init.ar"), filepath.Join(ex, "modules", path), filepath.Join(executable, "modules", realpath, "init.ar"), filepath.Join(executable, "modules", path)}
	var p string
	for _, p = range pathsToTest {
		if FileExists(p) {
			break
		}
	}
	if importing[p] == true {
		return nil, fmt.Sprint("Recursive import of ", p, " detected.\n\nHow to fix?\nCurrently this error cannot be pin pointed by Argons error detection, but this can be caused by a circular import or by a module that imports itself. Check .ar files inside ", origin, " for circular imports, or maybe this is an issue with one of the libaries used in the project so maybe get in contact with the library author.")
	} else if importing[p] == nil {
		importing[p] = true
		modules[p] = map[string]variableValue{}
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		stringdata := string(data)
		moduledata := modules[p]
		if main {
			vars["__main__"] = variableValue{
				TYPE:   "const",
				VAL:    p,
				EXISTS: true,
				origin: p,
				FUNC:   false,
			}
		}
		_, e := runStr(stringdata, p, moduledata)
		if e != nil {
			return nil, e
		}
		importing[p] = false
	}
	return modules[p], nil
}

func runStr(str string, origin string, variables map[string]variableValue) ([][]any, any) {
	translated, err := translate(str)
	if err != nil {
		return nil, err
	}
	fileparams := map[string]variableValue{"__file__": {
		TYPE:   "const",
		VAL:    origin,
		EXISTS: true,
		origin: origin,
		FUNC:   false,
	}}
	ty, _, resp := run(translated, origin, []map[string]variableValue{vars, fileparams, variables})
	if ty != nil {
		if ty == "error" {
			return nil, resp[len(resp)-1][0]
		}
		return nil, fmt.Sprint(ty) + " at top level"
	}
	return resp, nil
}
