package structs

import (
	"regexp"
	"time"
)

// file profile.yaml structure
type Config struct {
	Profiles []Profile `yaml:"profiles"`
}

type Profile struct {
	Pattern   string   `yaml:"p"`
	A         []string `yaml:"a"`
	R         []string `yaml:"r"`
	PatternRE *regexp.Regexp
	AcceptREs []*regexp.Regexp
	RejectREs []*regexp.Regexp
}

type Etalon struct {
	FileName string
	SHA1     string
	Date     time.Time
	Size     int64
}
