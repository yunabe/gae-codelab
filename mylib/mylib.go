// Package mylib implements utilities.
package mylib

import (
	"fmt"
)

// GetHelloMessage returns a hello message.
func GetHelloMessage(name string) string {
	return fmt.Sprintf("Hello %s", name)
}
