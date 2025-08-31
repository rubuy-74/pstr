package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rubuy-74/pstr/internal/parser"
)

func getRegex() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter regex:")
	fmt.Print("> ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	return input
}

func main() {
	for {
		regexString := getRegex()
		ctx, err := parser.Parse(regexString)
		if err != nil {
			log.Fatalf("Error while parsing regex: %s", err)
		}
		ctx.Print()
	}
}
