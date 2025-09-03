package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	StartOfText uint8 = 2
	EndOfText   uint8 = 3
)

const Infinite = -1
const Epsilon uint8 = 0

func GetChar(input string, pos int) uint8 {
	if pos >= len(input) {
		return EndOfText
	}

	if pos < 0 {
		return StartOfText
	}

	return uint8(input[pos])
}

func GetInput(message string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(message)
	fmt.Print("> ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	return input
}
