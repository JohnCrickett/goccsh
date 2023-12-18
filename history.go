package main

import (
	"bufio"
	"os"
	"path/filepath"
)

type History struct {
	commands []string
	file     *os.File
}

func NewHistory() *History {
	path, _ := os.UserHomeDir()
	var history_filename = filepath.Join(path, ".ccsh_history")
	file, _ := os.OpenFile(history_filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var commands []string

	for scanner.Scan() {
		commands = append(commands, scanner.Text())
	}

	return &History{commands: commands, file: file}
}

func shouldNotRecordInHistory(cmd string) bool {
	switch cmd {
	case
		"exit",
		"history":
		return true
	}
	return false
}

func (h *History) AddCommand(command string) {
	if shouldNotRecordInHistory(command) {
		return
	}
	h.commands = append(h.commands, command)
	h.file.WriteString(command + "\n")
}
