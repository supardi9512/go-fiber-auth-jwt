package helpers

import (
	"fmt"
	"net/mail"
)

func ValidateRegisterName(name string) error {
	if name == "" {
		return fmt.Errorf("Name is required")
	}

	if len(name) < 2 {
		return fmt.Errorf("Name must contain a minimum of 2 characters")
	}

	if len(name) > 50 {
		return fmt.Errorf("Name must contain a maximum of 50 characters")
	}

	return nil
}

func ValidateRegisterEmail(email string, compareEmails ...string) error {
	// Default compareEmail to an empty string
	compareEmail := ""
	if len(compareEmails) > 0 {
		compareEmail = compareEmails[0]
	}

	if email == "" {
		return fmt.Errorf("Email is required")
	}

	if len(email) > 50 {
		return fmt.Errorf("Email must contain a maximum of 50 characters")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("Incorrect email format")
	}

	if compareEmail != "" && compareEmail == email {
		return fmt.Errorf("Email already exists")
	}

	return nil
}

func ValidateRegisterUsername(username string, compareUsernames ...string) error {
	// Default compareUsername to an empty string
	compareUsername := ""
	if len(compareUsernames) > 0 {
		compareUsername = compareUsernames[0]
	}

	if username == "" {
		return fmt.Errorf("Username is required")
	}

	if len(username) < 5 {
		return fmt.Errorf("Username must contain a minimum of 5 characters")
	}

	if len(username) > 100 {
		return fmt.Errorf("Username must contain a maximum of 100 characters")
	}

	if compareUsername != "" && compareUsername == username {
		return fmt.Errorf("Username already exists")
	}

	return nil
}

func ValidateRegisterPassword(password string) error {
	if password == "" {
		return fmt.Errorf("Password is required")
	}

	if len(password) < 5 {
		return fmt.Errorf("Password must contain a minimum of 5 characters")
	}

	if len(password) > 100 {
		return fmt.Errorf("Password must contain a maximum of 100 characters")
	}

	return nil
}

func ValidateRegisterConfirmPassword(password string, confirmPassword string) error {
	if confirmPassword == "" {
		return fmt.Errorf("Confirm Password is required")
	}

	if password != confirmPassword {
		return fmt.Errorf("Confirm Password is not the same as Password")
	}

	return nil
}
