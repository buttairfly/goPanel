package pixelpipe

import "fmt"

// OutputIDMismatchError is an error which occours when a pipe output id could not be found
func OutputIDMismatchError(origin string, id, lookupID ID) error {
	return fmt.Errorf("could not find output with id %s of %s %s", id, origin, lookupID)
}
