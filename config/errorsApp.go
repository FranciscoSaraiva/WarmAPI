package config

import "fmt"

/*
*	Public Functions
 */

// CheckError function that receives an error and checks if it's not nil.
// If it's not nil it'll print the error to the console and return a true value.
// Else it returns a false value.
func CheckError(err error) bool {
	if err != nil {
		fmt.Printf("\nError: %s\n", err)
		return true
	}
	return false
}
