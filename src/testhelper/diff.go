package testhelper

import (
	"os/exec"
	"testing"
)

func Diff(t *testing.T, expectedFile, actualFile string) {
	cmd := []string{"diff", "-cbB", expectedFile, actualFile}
	RunCmd(t, cmd)
}

func RunCmd(t *testing.T, args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	t.Logf("RunCmd: %v", cmd.Args)
	output, err := cmd.CombinedOutput()
	message := string(output)
	if err != nil {
		t.Fatalf("err: %v cmd: %v msg:\n%v", err, cmd, message)
	}
	/*
		if message != "" {
			t.Log(message)
		}*/
}
