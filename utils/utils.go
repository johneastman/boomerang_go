package utils

import "fmt"

// TODO: move to utils package. Temporarily here to avoid import loop
func CreateError(err error) error {
	return fmt.Errorf(err.Error())
}
