package main

import (
	"os"
	"os/exec"
)

func updateEnvironment(env Environment) {
	if env == nil {
		return
	}
	for name, value := range env {
		if len(value) > 0 {
			os.Setenv(name, value)
		} else {
			os.Unsetenv(name)
		}
	}
}

// RunCmd runs a command + arguments (cmd) with environment variables from env
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 0
	}
	run := cmd[0]
	cmd = cmd[1:]
	updateEnvironment(env)
	command := exec.Command(run, cmd...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
	}
	return 0
}
