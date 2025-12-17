//go:build !windows && !linux
// +build !windows,!linux

package utils

import (
	"fmt"
	structs "zefc/structs"
)

func PrintGui(showGUI bool, msg string) {
	fmt.Print(msg)
	// No GUI support on this platform
}

func ShowGUI(zipFile string, errors int, zipCount int, pattern structs.Profile) {
	// No GUI support on this platform
	// Results are already printed to stdout
}
