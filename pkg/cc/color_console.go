package cc

import (
	"fmt"
	"time"
)

const (
	INFO    = "\033[96m"
	OKGREEN = "\033[92m"
	WARNING = "\033[93m"
	ERROR   = "\033[91m"
	BOLD    = "\033[1m"
	HEADING = "\x1b[44m %s \x1b[0m"
	ENDC    = "\033[0m"
)

func Warn(text string) {
	print(WARNING, text)
}

func Err(text string) {
	print(ERROR, text)
}

func Ok(text string) {
	print(OKGREEN, text)
}

func Info(text string) {
	print(INFO, text)
}

func Heading(text string) {
	fmt.Printf("%s%s", fmt.Sprintf(HEADING, text), ENDC)
}

func Bold(text string) {
	print(BOLD, text)
}

func print(style string, text string) {
	fmt.Printf("[SM] %s - %s%s%s\n", time.Now().Format("2006-01-02 15:04:05"), style, text, ENDC)
}
