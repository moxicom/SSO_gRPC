package utils

import "fmt"

func MakeError(operation string, error error) error {
	return fmt.Errorf("%s: %w", operation, error)
}
