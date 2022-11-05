package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func CreateError(lineNum int, errorMessage string, args ...any) error {
	errorMessagePrefix := fmt.Sprintf("error at line %d", lineNum)

	formattedErrorMessage := fmt.Sprintf(errorMessage, args...)
	fullErrorMessage := fmt.Sprintf("%s: %s", errorMessagePrefix, formattedErrorMessage)

	return fmt.Errorf(fullErrorMessage)
}

func NotANumberError(lineNum int, value string) error {
	return CreateError(lineNum, "cannot convert %#v to a number", value)
}

func CheckTypeError(lineNum int, actualType string, expectedType string) error {
	if expectedType != actualType {
		return CreateError(lineNum, "expected %s, got %s", expectedType, actualType)
	}
	return nil
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

func ConvertStringToInteger(value string) *int {
	integer, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return &integer
}

func ConvertStringToFloat(value string) *float64 {
	float, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil
	}
	return &float
}

func FloatToString(value float64) string {
	return fmt.Sprint(value)
}

func IntToString(value int) string {
	return fmt.Sprint(value)
}

func UserInput(prompt string) string {
	fmt.Printf("%s: ", prompt)

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line := scanner.Text()
		return line
	}
	return ""
}

func GetSource(path string) string {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(fileContent)
}
