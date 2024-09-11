package errors

import "fmt"

func StringLengthErrorMessage(min, max int) string {
	return fmt.Sprintf("Length must be between %d and %d characters.", min, max)
}
