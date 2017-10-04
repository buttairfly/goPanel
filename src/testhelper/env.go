package testhelper

import (
	"os"
	"strconv"
)

func RecordCall() bool {
	param := os.Getenv("TEST_RECORD")
	record, err := strconv.ParseBool(param)
	if err != nil {
		return false
	}
	return record
}
