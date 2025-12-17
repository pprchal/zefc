package utils

import (
	"os"
	"path/filepath"
	"regexp"
	structs "zefc/structs"

	"github.com/goccy/go-yaml"
)

// load config from $HOME/zefc/patterns.yaml
func LoadConfig() structs.Config {
	path := filepath.Join(GetHomeDir(), "zefc", "patterns.yaml")
	// check if file exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		panic("Config file not found: " + path)
	}

	yml, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var config structs.Config
	if err := yaml.Unmarshal([]byte(yml), &config); err != nil {
		panic(err)
	}

	for i := range config.Profiles {
		CompileRegexps(&config.Profiles[i])
	}
	return config
}

// prepare compiled regexps in profile
func CompileRegexps(profile *structs.Profile) {
	profile.PatternRE = regexp.MustCompile(profile.Pattern)
	for i := range profile.R {
		reject := regexp.MustCompile(profile.R[i])
		profile.RejectREs = append(profile.RejectREs, reject)
	}

	for i := range profile.A {
		accept := regexp.MustCompile(profile.A[i])
		profile.AcceptREs = append(profile.AcceptREs, accept)
	}
}
