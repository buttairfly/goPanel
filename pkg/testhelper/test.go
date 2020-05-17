package testhelper

import (
	"os"
	"testing"
)

// FailAndSkip fails the test but does not abort other tests
func FailAndSkip(t *testing.T, args ...interface{}) {
	t.Log(args...)
	t.Fail()
	t.SkipNow()
}

// FileExistsOrSkip tests if a file at fullPath is available or skips the test
func FileExistsOrSkip(t *testing.T, fullPath string) {
	if !RecordCall() {
		if _, err := os.Stat(fullPath); err != nil {
			t.Log(err.Error())
			FailAndSkip(t, "Re-Run: env TEST_RECORD=true go test ./...")
		}
	}
}
