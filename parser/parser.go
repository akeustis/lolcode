package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Parse echoes each line of stdin back out
func Parse() {
	stdin := bufio.NewReader(os.Stdin)
	for {
		text, err := stdin.ReadString('\n')
		if err != nil {
			return
		}
		tokens := strings.Fields(text)
		for _, token := range tokens {
			fmt.Printf("%s,", token)
		}
		fmt.Println()
	}
}
