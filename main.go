package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("ccsh> ")
	input, _ := reader.ReadString('\n')

	input = strings.TrimSpace(input)
	parts := strings.Fields(input)

	var command = parts[0]
	var args = parts[1:]

	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
