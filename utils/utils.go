package utils

import (
	"fmt"
)

func CreateError(lineNum int, errorMessage string, args ...any) error {
	errorMessagePrefix := fmt.Sprintf("error at line %d", lineNum)

	formattedErrorMessage := fmt.Sprintf(errorMessage, args...)
	fullErrorMessage := fmt.Sprintf("%s: %s", errorMessagePrefix, formattedErrorMessage)

	return fmt.Errorf(fullErrorMessage)
}
