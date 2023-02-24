package util

import (
	"os"
	"os/exec"
)

func ExecBashCmd(sh string) error {
	cmd := exec.Command("bash", "-c", sh)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ExecZshCmd(sh string) error {
	cmd := exec.Command("zsh", "-c", sh)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
