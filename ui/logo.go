package ui

import (
	"fmt"
	"strings"
)

func PrintLogo() {
	fmt.Println(" _   _      _       _ _          ____       _          _")
	fmt.Println("| | | | ___| |_ __ | (_)_ __    / ___| ___ | |__   ___| |_ __   ___ _ __")
	fmt.Println("| |_| |/ _ \\ | '_ \\| | | '_ \\  | |  _ / _ \\| '_ \\ / _ \\ | '_ \\ / _ \\ '__|")
	fmt.Println("|  _  |  __/ | |_) | | | |_) | | |_| | (_) | | | |  __/ | |_) |  __/ |")
	fmt.Println("|_| |_|\\___|_| .__/|_|_| .__/   \\____|\\___/|_| |_|\\___|_| .__/ \\___|_|")
	fmt.Println("             |_|       |_|                               |_|")
	fmt.Println()
}

func PromptYesNo(prompt string, defaultYes bool) bool {
	def := "y"
	if !defaultYes {
		def = "n"
	}
	v := Prompt(prompt, def)
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "y", "yes", "1":
		return true
	case "n", "no", "0":
		return false
	default:
		return defaultYes
	}
}

