package main

import (
	"fmt"
	"github.com/MichalMichalak/arx2/conf"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	r := conf.NewResolver([]string{"_temp/config.yaml"})
	c := r.Conf()
	for k, v := range c {
		fmt.Printf("conf[%s]=`%v`\n", k, v)
	}
	for _, w := range r.Warns() {
		fmt.Println("!!! WARN: " + w)
	}
	fmt.Println()
	inject(r)
}

func inject(r conf.Resolver) {
	type host struct {
		Name string `conf:"service.host"`
		Port int    `conf:"service.port"`
	}
	type configuration struct {
		RootSecondLevel bool   `conf:"root.second-level"`
		F               string `conf:"f"`
		Pwd             string `conf:"config.path"`
		Host            host   `conf:""`
		NonConf         int64
		TagButNonConf   []string `some:"array"`
	}
	c := configuration{}
	err := conf.Configure(&c, r)
	if err != nil {
		panic(err)
	}
	fmt.Println("-----")
	spew.Dump(c)
}
