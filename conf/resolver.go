package conf

import (
	"github.com/unravelin/core/arx2/log"
	"strings"
)

type DefaultResolver struct {
	logger  log.Logger
	loaders []genericLoader
}

// NewResolver creates new resolver that will load config values from shell environment, YAML files and finally from
// command line arguments.
//
// Source priority, starting from the lowest: Shell environment; YAML files, in order provided; Command line arguments.
func NewResolver(logger log.Logger, paths []string) DefaultResolver {
	loaders := []genericLoader{newEnvLoader()}
	for _, p := range paths {
		loaders = append(loaders, newYamlLoader(p))
	}
	loaders = append(loaders, newCmdLoader())
	for _, l := range loaders {
		l.load()
	}
	return DefaultResolver{logger: logger, loaders: loaders}
}

func (r DefaultResolver) Conf() map[string]interface{} {
	c := map[string]interface{}{}
	for _, l := range r.loaders {
		for k, v := range l.conf() {
			c[k] = v
		}
	}
	return c
}

func (r DefaultResolver) Warns() []string {
	var warns []string
	for _, l := range r.loaders {
		if len(l.warns()) == 0 {
			continue
		}
		s := l.name() + " warnings: " + strings.Join(l.warns(), "; ")
		warns = append(warns, s)
	}
	return warns
}
