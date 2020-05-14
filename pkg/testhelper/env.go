package testhelper

import (
	"os"
	"strconv"
)

// RecordCall checks environment variables to check if a new test record should be taken before the test
func RecordCall() bool {
	param := os.Getenv("TEST_RECORD")
	isRecordCall, err := strconv.ParseBool(param)
	if err != nil {
		return false
	}
	return isRecordCall
}
