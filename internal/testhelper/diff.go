package testhelper

import (
	"os/exec"
	"testing"
)

// Diff runs the diff program with two file paths
func Diff(t *testing.T, expectedFile, actualFile string) {
	cmd := []string{"diff", "-abBcd", expectedFile, actualFile}
	RunCmd(t, cmd)
}

// RunCmd runs an linux program for a test with args
func RunCmd(t *testing.T, args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	t.Logf("RunCmd: %v", cmd.Args)
	output, err := cmd.CombinedOutput()
	message := string(output)
	if err != nil || message != "" {
		t.Fatalf("err: %v", message)
	}
}
