package conf

import (
	"os"
	"regexp"
)

type cmdLoader struct {
	regex    *regexp.Regexp
	warnings []string
	values   map[string]interface{}
}

func newCmdLoader() *cmdLoader {
	regex := regexp.MustCompile("--([a-zA-Z][-.a-zA-Z0-9]*)=(.+)")
	val := map[string]interface{}{}
	return &cmdLoader{regex: regex, values: val}

}

func (l *cmdLoader) load() {
	for _, arg := range os.Args {
		found := l.regex.FindStringSubmatch(arg)
		if found != nil {
			key := normalizeCmdLineArgKey(found[1])
			l.addValue(key, found[2])
		}
	}
}

func (l *cmdLoader) addValue(key string, value string) {
	if key == "" {
		return
	}
	if _, exist := l.values[key]; exist {
		warn := "key `" + key + "` already exists, value `" + value + "` ignored"
		l.warnings = append(l.warnings, warn)
		return
	}
	l.values[key] = value
	return
}

const cmdLoaderName = "cmd"

func (l *cmdLoader) name() string {
	return cmdLoaderName
}

func (l *cmdLoader) conf() map[string]interface{} {
	return l.values
}

func (l *cmdLoader) warns() []string {
	return l.warnings
}
