package main

import (
	"fmt"
	"os"
)

func shell() {
	fmt.Println("Argon Shell\nMIT License\n\nCopyright (c) 2022 William Bell\ncall 'license' for license information\n ")
	indent := 0
	temp := ""
	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	variables := map[string]variableValue{}
	for {
		tempstr := ""
		if indent >= 0 {
			val, _ := ArgonInput("> ")
			tempstr = val.(string)
		} else {
			val, _ := ArgonInput("... ")
			tempstr = val.(string)
		}
		if opensIndent(tempstr) {
			indent--
		} else if closeCompile.MatchString(tempstr) {
			indent++
		}
		temp += tempstr + "\n"
		if indent >= 0 {
			resp, err := runStr(temp, ex, variables)
			if err != nil {
				fmt.Println(err)
			}
			for i := 0; i < len(resp); i++ {
				if resp[i][1] != nil {
					fmt.Println(anyToArgon(resp[i][0], false))
				}
			}
			temp = ""
		}
	}
}
