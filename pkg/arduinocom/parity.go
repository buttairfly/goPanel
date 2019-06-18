package arduinocom

import (
	"fmt"
	"strings"
)

// CalcHexParity takes a bare string command, calculated the parity and returns the full string with return line
func CalcHexParity(bareCommand string, seed byte) string {
	parity := calcHexParityChar(bareCommand, seed)
	return fmt.Sprintf("%s%s\n", bareCommand, parity)
}

// CheckHexParity checks the correctness of the parity char
func CheckHexParity(checkCommand string, seed byte) bool {
	checkCommand = strings.TrimRight(checkCommand, "\n")
	if len(checkCommand) < 1 {
		return false
	}
	parityPos := len(checkCommand) - 1
	receivedParity := checkCommand[parityPos]
	bareCommand := checkCommand[0:parityPos]
	calcParity := calcHexParityChar(bareCommand, seed)
	if calcParity != string(receivedParity) {
		return false
	}
	return true
}

func calcHexParityChar(bareCommand string, seed byte) string {
	parity := seed
	for i := 0; i < len(bareCommand); i++ {
		parity ^= bareCommand[i]
	}
	highParity := (parity >> 4) & 0xf
	lowParity := parity & 0xf
	parity = highParity ^ lowParity
	return fmt.Sprintf("%01x", parity)
}
