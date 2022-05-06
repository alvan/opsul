package utils

import (
	"errors"
	"fmt"
)

func Args(text string) ([]string, error) {
	args := make([]string, 0)

	state := ""
	value := ""
	quote := ""
	slash := false

	for i := 0; i < len(text); i++ {
		char := text[i]

		if slash {
			value += string(char)
			slash = false
			continue
		}

		if char == '\\' {
			slash = true
			continue
		}

		if state == "quote" {
			if string(char) == quote {
				state = ""
				quote = ""
				if value != "" {
					args = append(args, value)
					value = ""
				}
			} else {
				value += string(char)
			}

			continue
		}

		if char == '"' || char == '\'' {
			state = "quote"
			quote = string(char)

			if value != "" {
				args = append(args, value)
				value = ""
			}

			continue
		}

		if char == ' ' || char == '\t' {
			if value != "" {
				args = append(args, value)
				value = ""
			}

			continue
		}

		value += string(char)
	}

	if state == "quote" || slash {
		return []string{}, errors.New(fmt.Sprintf("Invalid args in command line: %s", text))
	}

	if value != "" {
		args = append(args, value)
	}

	return args, nil
}
