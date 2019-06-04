package testhelper

import (
	"os"
	"strconv"
)

// RecordCall checks environment variables to
func RecordCall() bool {
	param := os.Getenv("TEST_RECORD")
	record, err := strconv.ParseBool(param)
	if err != nil {
		return false
	}
	return record
}
