package utils

import (
	"fmt"
	"strconv"
)

func CreateError(lineNum int, errorMessage string, args ...any) error {
	errorMessagePrefix := fmt.Sprintf("error at line %d", lineNum)

	formattedErrorMessage := fmt.Sprintf(errorMessage, args...)
	fullErrorMessage := fmt.Sprintf("%s: %s", errorMessagePrefix, formattedErrorMessage)

	return fmt.Errorf(fullErrorMessage)
}

func ConvertStringToInteger(lineNum int, value string) (*int, error) {
	integer, err := strconv.Atoi(value)
	if err != nil {
		return nil, CreateError(lineNum, "list index must be an integer")
	}
	return &integer, nil
}

func CheckOutOfRange(lineNum int, index int, listLen int) error {
	if index < 0 || index > listLen-1 {
		return CreateError(
			lineNum,
			"index of %d out of range (%d to %d)",
			index, 0, listLen-1,
		)
	}
	return nil
}
