package main

import (
	"fmt"
	"log"

	"github.com/rubuy-74/pstr/internal/parser"
	"github.com/rubuy-74/pstr/internal/state_machine"
	"github.com/rubuy-74/pstr/internal/utils"
)

func MainDev() {
	for {
		regexString := utils.GetInput("Enter regex")
		stringToCheck := utils.GetInput("Enter string to check")

		ctx, err := parser.Parse(regexString)
		if err != nil {
			log.Fatalf("Error while parsing regex: %s", err)
		}
		ctx.Print()

		nfa, err := state_machine.ToNFA(ctx)
		if err != nil {
			log.Fatalf("Error while creating NFA: %s", err)
		}
		fmt.Println(nfa)

		isValid := nfa.Check(stringToCheck, -1)
		if isValid {
			fmt.Println("Congratulations, the string is VALID")
		} else {
			fmt.Println("Fuck you, the string is NOT VALID")
		}
	}
}
