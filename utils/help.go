package utils

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"runtime"
)

func GetHomeDir() string {
	if currentUser, err := user.Current(); err == nil && currentUser.HomeDir != "" {
		return currentUser.HomeDir
	}

	// Fallback to environment variables
	switch runtime.GOOS {
	case "windows":
		return os.Getenv("USERPROFILE")
	default:
		return os.Getenv("HOME")
	}
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func PrintHelp() {
	fmt.Println("zip content verifier, version 1.2")
	fmt.Println("zefc <zip_file> [--gui]")
	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("")
	fmt.Println("If no etalon exists for the module, it will be created.")
	fmt.Println("Working directory for etalon files is: %USERPROFILE%\\zefc")
	fmt.Println("If --gui is provided, a message box will be shown with the results.")
	fmt.Println("")
	fmt.Println(`1. resolve patterns from %USERPROFILE%\zefc\patterns => <module>`)
	fmt.Println(`2. load etalon from %USERPROFILE%\zefc\<module>.eta`)
	fmt.Println(`3. calculate sha1 hashes of files in the zip and compare with etalon`)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Note: if you want to build new etalon, just delete the existing ETA_FILE.")
	fmt.Println("Pavel Prchal, prchalp@gmail.com")
}
