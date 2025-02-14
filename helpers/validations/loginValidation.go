package helpers

import (
	"fmt"
)

func ValidateLoginUsername(username string) error {
	if username == "" {
		return fmt.Errorf("User is required")
	}

	return nil
}

func ValidateLoginPassword(password string) error {
	if password == "" {
		return fmt.Errorf("Password is required")
	}

	return nil
}
