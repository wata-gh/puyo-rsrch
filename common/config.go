package puyorsrch

import "github.com/BurntSushi/toml"

type Config struct {
	PuyoConfig PuyoConfig
}

type PuyoConfig struct {
	SearchLocations [][2]int
	Field           string
	Mappings        map[string]int
	ExpectedChains  int
	AllClear        bool
	ExportImagePath string
}

func ReadConfig() *Config {
	var conf Config
	toml.DecodeFile("puyo.toml", &conf)
	if conf.PuyoConfig.SearchLocations == nil {
		panic("conf.PuyoConfig.SearchLocations must be set.")
	}
	if conf.PuyoConfig.Mappings == nil {
		conf.PuyoConfig.Mappings = map[string]int{"a": 4, "b": 5, "c": 6, "d": 7, " ": 0}
	}
	return &conf
}
