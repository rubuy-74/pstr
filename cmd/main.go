package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rubuy-74/pstr/internal/parser"
	"github.com/rubuy-74/pstr/internal/state_machine"
)

func getInput(message string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(message)
	fmt.Print("> ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	return input
}

func main() {
	for {
		regexString := getInput("Enter regex")
		stringToCheck := getInput("Enter string to check")

		ctx, err := parser.Parse(regexString)
		if err != nil {
			log.Fatalf("Error while parsing regex: %s", err)
		}
		ctx.Print()

		nfa := state_machine.ToNFA(ctx)
		fmt.Println(nfa)

		isValid := nfa.Check(stringToCheck, -1)
		if isValid {
			fmt.Println("Congratulations, the string is VALID")
		} else {
			fmt.Println("Fuck you, the string is NOT VALID")
		}
	}
}
