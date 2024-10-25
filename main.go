package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func notifyCmds(cmds []*exec.Cmd, s os.Signal) {
	for _, cmd := range cmds {
		cmd.Process.Signal(s)
	}
}

func main() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	var cmds []*exec.Cmd
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		select {
		case sig := <-signalChannel:
			switch sig {
			case os.Interrupt, syscall.SIGTERM:
				notifyCmds(cmds, sig)
			}
		case <-ctx.Done():
			close(signalChannel)
			return
		}

	}(ctx)

	history := NewHistory()

	for {
		fmt.Print("ccsh> ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		commands := strings.Split(input, "|")
		var output io.ReadCloser

		for _, commandLine := range commands {
			commandLine = strings.TrimSpace(commandLine)
			history.AddCommand(commandLine)

			parts := strings.Fields(commandLine)
			var command = parts[0]
			var args = parts[1:]

			switch command {
			case "exit":
				cancel()
				history.file.Close()
				os.Exit(0)

			case "cd":
				var path string
				var err error

				if len(args) > 0 {
					path = args[0]
				} else {
					path, _ = os.UserHomeDir()
				}
				err = os.Chdir(path)
				if err != nil {
					fmt.Printf("%v\n", err)
				}

			case "history":
				for i, hc := range history.commands {
					fmt.Printf("%d: %s\n", i, hc)
				}

			case "pwd":
				pr, pw := io.Pipe()
				go Pwd(pw)
				if len(commands) == 1 {
					io.Copy(os.Stdout, pr)
				} else {
					output = pr
				}

			default:
				cmd := exec.Command(command, args...)
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr

				cmds = append(cmds, cmd)

				if output != nil {
					cmd.Stdin = output
				}
				output, _ = cmd.StdoutPipe()
			}
		}

		if len(cmds) > 0 {
			cmds[len(cmds)-1].Stdout = os.Stdout
		}

		for _, cmd := range cmds {
			cmd.Start()
		}

		for _, cmd := range cmds {
			err := cmd.Wait()
			if err != nil {
				if err.Error() != "signal: interrupt" {
					if cmd.ProcessState.ExitCode() == -1 {
						fmt.Printf("command not found: %s\n", cmd.Path)
					}
				}
			}
		}
		cmds = nil
	}
}
