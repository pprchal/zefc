//go:build windows
// +build windows

package utils

import (
	"fmt"
	structs "zefc/structs"

	"github.com/zzl/go-win32api/v2/win32"
)

func PrintGui(showGUI bool, msg string) {
	fmt.Print(msg)
	if showGUI {
		win32.MessageBoxW(0, win32.StrToPwstr(msg), win32.StrToPwstr("zefc"), win32.MB_OK)
	}
}

func ShowGUI(zipFile string, errors int, zipCount int, pattern structs.Profile) {
	if errors == 0 {
		win32.MessageBoxW(0, win32.StrToPwstr(fmt.Sprintf("All %d files are ok", zipCount)), win32.StrToPwstr(zipFile), win32.MB_OK)
		return
	}

	msg := fmt.Sprintf("Some files are missing or differ - %d !\n", errors)
	win32.MessageBoxW(0, win32.StrToPwstr(msg), win32.StrToPwstr(zipFile), win32.MB_OK|win32.MB_ICONEXCLAMATION)
}
