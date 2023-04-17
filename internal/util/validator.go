package util

import (
	"errors"
	"fmt"
	"net/mail"
	"os"
	"path/filepath"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("must contain only lowercase letters, digits, or underscore")
	}
	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidFullName(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}
	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 5, 100)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not a valid email address")
	}
	return nil
}

func ValidateBalance(value int64) error {
	if value <= 0 {
		return fmt.Errorf("must be greater than 0")
	}
	return nil
}

func ValidateFolderAndFile(folder, file string) error {
	// Get the absolute path of the folder/file from root project directory
	folderPath := filepath.Join(".", folder) // Use "." for the current directory
	filePath := filepath.Join(".", folder, file)
	// Check if a folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.Mkdir(folderPath, 0755)
		if err != nil {
			return errors.New(err.Error())
		}
	}

	// Check if a file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create the file with permission 0664
		file, err := os.Create(filePath)
		if err != nil {
			return errors.New(err.Error())
		}
		defer file.Close()

		// Set the file permission to 0664
		err = file.Chmod(0664)
		if err != nil {
			return errors.New(err.Error())
		}
	}

	return nil
}
