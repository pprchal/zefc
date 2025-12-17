package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	structs "zefc/structs"
)

// save etalon to $HOME/zefc/{module}.eta
func SaveEtalon(module string, etalon chan structs.Etalon, zipFile string) {
	path := filepath.Join(GetHomeDir(), "zefc", module+".eta")

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("# zefc etalon file\n")
	fmt.Fprintf(writer, "# zip source: %s\n", zipFile)

	for e := range etalon {
		line := strings.TrimSpace(e.SHA1) + " " + strings.TrimSpace(e.FileName) + " " + strconv.FormatInt(e.Size, 10) + " \n"
		_, err := writer.WriteString(line)
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()
}

// load etalon from $HOME/zefc/{module}.eta
func LoadEtalon(module string) ([]structs.Etalon, string) {
	var results []structs.Etalon

	path := filepath.Join(GetHomeDir(), "zefc", module+".eta")

	file, err := os.Open(path)
	if err != nil {
		return results, path
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		split := strings.SplitN(line, " ", 4)
		if len(split) == 4 {
			size, _ := strconv.ParseInt(strings.TrimSpace(split[2]), 10, 64)
			etalon := structs.Etalon{
				FileName: strings.TrimSpace(split[1]),
				SHA1:     strings.TrimSpace(split[0]),
				Size:     size,
			}
			results = append(results, etalon)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return results, path
}

// IsHandledFile checks if the file is a .zip and matches any of the loaded patterns.
func IsHandledFile(file string, config structs.Config) (bool, string, structs.Profile) {
	if !strings.HasSuffix(file, ".zip") {
		return false, "", structs.Profile{}
	}

	for _, profile := range config.Profiles {
		matches := profile.PatternRE.FindStringSubmatch(file)
		if len(matches) > 0 {
			return true, matches[1], profile
		}
	}

	return false, "", structs.Profile{}
}
