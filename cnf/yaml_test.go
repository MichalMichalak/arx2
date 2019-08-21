package cnf

import (
	"fmt"
	"strings"
	"testing"
)

func TestYaml(t *testing.T) {
	path := "_temp/config.yaml"
	l := newYamlLoader(path)
	l.load()
	println(strings.Join(l.warns(), ";"))
	for k, v := range l.values {
		fmt.Printf("conf[%s]=`%v`(%T)\n", k, v, v)
	}
}
