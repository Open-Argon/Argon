package main

import "fmt"

func shell() {
	fmt.Println("Argon Shell\nMIT License\n\nCopyright (c) 2022 William Bell\ncall 'license' for license information\n ")
	indent := 0
	temp := ""
	for {
		tempstr := ""
		if indent >= 0 {
			tempstr = ArgonInput("> ").(string) + "\n"
		} else {
			tempstr = ArgonInput("... ").(string) + "\n"
		}
		if !switchCloseCompile.MatchString(tempstr) {
			if openCompile.MatchString(tempstr) {
				indent--
			} else if closeCompile.MatchString(tempstr) {
				indent++
			}
		}
		temp += tempstr
		if indent >= 0 {
			resp := runStr(temp)
			for i := 0; i < len(resp); i++ {
				if resp[i][1] != nil {
					fmt.Println(anyToArgon(resp[i][0]))
				}
			}
			temp = ""
		}
	}
}
