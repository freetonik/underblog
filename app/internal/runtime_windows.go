// +build windows
package internal

const defaultWLimit = 100

// Return workers limit
// todo: test it on Шindows and set more effective limit
func GetWorkersLimit() int {
	return defaultWLimit
}
