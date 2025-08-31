package utils

const (
	startOfText uint8 = 2
	endOfText   uint8 = 3
)

const Infinite = -1
const Epsilon uint8 = 0

func GetChar(input string, pos int) uint8 {
	if pos >= len(input) {
		return endOfText
	}

	if pos < 0 {
		return startOfText
	}

	return uint8(input[pos])
}
