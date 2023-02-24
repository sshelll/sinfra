package util

import "testing"

func TestExecBashCmd(t *testing.T) {
	if err := ExecBashCmd("echo 'hello world'"); err != nil {
		t.Errorf("ExecBashCmd() error = %v", err)
	}
}

func TestExecZshCmd(t *testing.T) {
	if err := ExecZshCmd("echo 'hello world'"); err != nil {
		t.Errorf("ExecZshCmd() error = %v", err)
	}
}
