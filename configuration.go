package main

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	FirstBoot        bool   `json:"firstBoot"`
	OnFirstBoot      string `json:"onFirstBoot"`
	OnSubsequentBoot string `json:"onSubsequentBoot"`
}

func Load(path string) (*Configuration, error) {
	configuration := new(Configuration)
	configuration.init()

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return configuration, nil
	} else if err != nil {
		return nil, err
	}
	defer func() {
		err = file.Close()
	}()

	decoder := json.NewDecoder(file)

	if decoder.More() {
		if err = decoder.Decode(configuration); err != nil {
			return nil, err
		}
	}
	return configuration, nil
}

func (conf *Configuration) init() {
	conf.FirstBoot = false
	conf.OnFirstBoot = ""
	conf.OnSubsequentBoot = "/sbin/init"
}

func (conf Configuration) OnBoot() string {
	if conf.FirstBoot && conf.OnFirstBoot != "" {
		return conf.OnFirstBoot
	}
	return conf.OnSubsequentBoot
}

func (conf Configuration) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(conf)
}
