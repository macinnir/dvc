package utils

import (
	"github.com/Tkanos/gonfig"
)

// LoadConfig loads a config file
func LoadConfig(path string, config interface{}) {
	e := gonfig.GetConf(path, config)
	if e != nil {
		panic(e)
	}
}
