package utility

import "strings"

// Contains() - checks if a string exists in a string array
// -- pass in an array of strings and the string you're looking for in that array
// -- returns a true if that string is in anywhere in the arrar
// -- returns false if it's not found
func Contains(s []string, str string) bool {
	var bCheck bool
	for _, v := range s {
		if strings.Contains(v, str) {
			bCheck = true
		} else {
			bCheck = false
		}
	}
	return bCheck
}
