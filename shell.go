package main

import "fmt"

func shell() {
	fmt.Println("Argon Shell\nMIT License\n\nCopyright (c) 2022 William Bell\ncall 'license' for license information\n ")
	for {
		runStr(ArgonInput("> ").(string))
	}
}
