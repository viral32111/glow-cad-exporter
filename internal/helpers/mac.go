package helpers

import "strings"

func NormaliseMACAddress(input string) string {
	if len(input) != 12 || strings.Contains(input, ":") {
		return input
	}

	var output []string

	for index := 0; index < len(input); index += 2 {
		output = append(output, input[index:index+2])
	}

	return strings.Join(output, ":")
}
