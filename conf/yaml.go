package conf

import (
	"gopkg.in/yaml.v2"
	"os"
	"strconv"
)

type yamlLoader struct {
	fileName string
	warnings []string
	values   map[string]interface{}
}

func newYamlLoader(path string) *yamlLoader {
	return &yamlLoader{fileName: path, warnings: []string{}, values: map[string]interface{}{}}
}

const yamlLoaderName = "yaml"

func (yamlLoader) name() string {
	return yamlLoaderName
}

func (l *yamlLoader) load() {
	f, err := os.Open(l.fileName)
	if err != nil {
		w := "failed to open file `" + l.fileName + "`: " + err.Error()
		l.warnings = append(l.warnings, w)
		return
	}
	defer func() {
		err := f.Close()
		if err != nil {
			w := "failed to close file `" + l.fileName + "`: " + err.Error()
			l.warnings = append(l.warnings, w)
		}
	}()
	dec := yaml.NewDecoder(f)
	raw := map[interface{}]interface{}{}
	err = dec.Decode(&raw)
	if err != nil {
		w := "failed to decode file`" + l.fileName + "`: " + err.Error()
		l.warnings = append(l.warnings, w)
		return
	}
	l.traverse(raw, "")
}

func (l yamlLoader) traverse(raw map[interface{}]interface{}, prefix string) {
	for k, v := range raw {
		addPrefix := ""
		if prefix != "" {
			addPrefix += "."
		}
		addPrefix += k.(string)
		switch v2 := v.(type) {
		case map[interface{}]interface{}:
			l.traverse(v2, prefix+addPrefix)
		case []interface{}:
			for i, elem := range v2 {
				finalPref := prefix + addPrefix + "." + strconv.Itoa(i)
				l.addValue(finalPref, elem)
			}
		default:
			l.addValue(prefix+addPrefix, v)
		}
	}
}

func (l yamlLoader) addValue(key string, value interface{}) {
	normalized := normalizeCmdLineArgKey(key)
	l.values[normalized] = value
}

func (l yamlLoader) conf() map[string]interface{} {
	return l.values
}

func (l yamlLoader) warns() []string {
	return l.warnings
}
