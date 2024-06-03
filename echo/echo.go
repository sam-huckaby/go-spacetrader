package echo

import (
	"errors"
	"fmt"
)

// Hello returns a greeting for the specified person
func Hello(name string) (string, error) {
	// If a name is not given, return an error
	if name == "" {
		return name, errors.New("No name specified")
	}

	return fmt.Sprintf("Hello %v, Welcome!", name), nil
}
