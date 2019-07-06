package conf

import (
	"os"
	"regexp"
)

type envLoader struct {
	regex    *regexp.Regexp
	warnings []string
	values   map[string]interface{}
	envVars  func() []string
}

func newEnvLoader() *envLoader {
	m := regexp.MustCompile("([a-zA-Z_][a-zA-Z0-9_]*)=(.+)")
	return &envLoader{regex: m, values: map[string]interface{}{}, envVars: os.Environ}
}

func (l *envLoader) load() {
	for _, v := range l.envVars() {
		found := l.regex.FindStringSubmatch(v)
		if found != nil {
			key := normalizeEnvVarKey(found[1])
			l.addValue(key, found[2])
		}
	}
}

func (l *envLoader) addValue(key string, value string) {
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

const envLoaderName = "env"

func (l *envLoader) name() string {
	return envLoaderName
}

func (l *envLoader) conf() map[string]interface{} {
	return l.values
}

func (l *envLoader) warns() []string {
	return l.warnings
}
