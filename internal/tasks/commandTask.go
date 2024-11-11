package tasks

import (
	"fmt"
	"os"
	"os/exec"
)

type CommandTask struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

func RunCommandTask(task CommandTask) (Register, error) {
	cmd := exec.Command(task.Command, task.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return Register{}, fmt.Errorf("error running command: %v", err)
	}

	return Register{}, nil
}
