package testhelper

import "testing"

// FailAndSkip fails the test but does not abort other tests
func FailAndSkip(t *testing.T, args ...interface{}) {
	t.Log(args...)
	t.Fail()
	t.SkipNow()
}
