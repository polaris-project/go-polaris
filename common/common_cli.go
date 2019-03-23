// Package common defines a set of commonly used helper methods and data types.
package common

import (
	"errors"
	"strings"
)

var (
	// ErrNilInput is an error definition describing input of 0 char length.
	ErrNilInput = errors.New("nil input")
)

/* BEGIN EXPORTED METHODS */

// ParseStringMethodCall attempts to parse string as method call, returning receiver, method name and params.
func ParseStringMethodCall(input string) (string, string, []string, error) {
	if input == "" { // Check for errors
		return "", "", []string{}, ErrNilInput // Return found error
	} else if !strings.Contains(input, "(") || !strings.Contains(input, ")") {
		input = input + "()" // Fetch receiver methods
	}

	if !strings.Contains(input, ".") { // Check for nil receiver
		return "", "", []string{}, errors.New("invalid method " + input) // Return found error
	}

	method := strings.Split(strings.Split(input, "(")[0], ".")[1] // Fetch method

	receiver := StringFetchCallReceiver(input) // Fetch receiver

	params := []string{} // Init buffer

	if strings.Contains(input, ",") || !strings.Contains(input, "()") { // Check for nil params
		params, _ = ParseStringParams(input) // Fetch params
	}

	return receiver, method, params, nil // No error occurred, return parsed method+params
}

// ParseStringParams attempts to fetch string parameters from (..., ..., ...) style call.
func ParseStringParams(input string) ([]string, error) {
	if input == "" { // Check for errors
		return []string{}, ErrNilInput // Return found error
	}

	parenthesesStripped := StringStripReceiverCall(input) // Strip parentheses

	params := strings.Split(parenthesesStripped, ", ") // Split by ', '

	return params, nil // No error occurred, return split params
}

// StringStripReceiverCall strips receiver from string method call.
func StringStripReceiverCall(input string) string {
	openParenthIndex := strings.Index(input, "(")      // Get open parent index
	closeParenthIndex := strings.LastIndex(input, ")") // Get close parent index

	return input[openParenthIndex+1 : closeParenthIndex] // Strip receiver
}

// StringFetchCallReceiver attempts to fetch receiver from string, as if it were an x.y(..., ..., ...) style method call.
func StringFetchCallReceiver(input string) string {
	return strings.Split(strings.Split(input, "(")[0], ".")[0] // Return split string
}

/* END EXPORTED METHODS */
