package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Prompt(prompt string, defaultValue string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	inputStr, _ := reader.ReadString('\n')
	inputStr = strings.TrimSpace(inputStr)
	if inputStr == "" {
		return defaultValue
	}
	return inputStr
}

func PromptRequired(prompt string) string {
	for {
		v := Prompt(prompt, "")
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
}
