package greetings

import (
	"errors"
	"fmt"
)

func Hello(name string) (string, error) {

	if name == "" {
		return "", errors.New("name can't be empty")
	}

	message := fmt.Sprintf("Hi, %v. Welcome to GoLang", name)
	return message, nil
}
