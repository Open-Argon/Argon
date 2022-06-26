package main

import "log"

var catch = [](func(any)){}

func err(err string) {
	if len(catch) > 0 {
		catch[len(catch)-1](err)
		catch[len(catch)-1] = nil
	} else {
		log.Fatal(err)
	}
}

func catchError(f func(any)) {
	catch = append(catch, f)
}
