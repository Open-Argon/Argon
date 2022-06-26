package main

import (
	"log"
	"os"
	"regexp"
	"strings"
)

type code struct {
	code string
	line int
}

func stringencode(str string) string {
	output := strings.Trim(str, " ")
	output = output[1 : len(output)-1]
	return output
}

func makeRegex(str string) *regexp.Regexp {
	Compile, err := regexp.Compile("^(" + str + ")$")
	if err != nil {
		log.Fatal(err)
	}
	return Compile
}

func removeNils(initialMap map[string]interface{}) map[string]interface{} {
    withoutNils := map[string]interface{}{}
    for key, value := range initialMap {
        _, ok := value.(map[string]interface{})
        if ok {
            value = removeNils(value.(map[string]interface{}))
            withoutNils[key] = value
            continue
        }
        if value != nil {
            withoutNils[key] = value
        }
    }
    return withoutNils
}

func main() {
	if len(os.Args) == 1 {
		shell()
	} else {
		importMod(os.Args[1])
	}
}
