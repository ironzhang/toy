package jsoncfg

import (
	"encoding/json"
	"io/ioutil"
)

func LoadFromFile(filename string, cfg interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, cfg); err != nil {
		return err
	}
	return nil
}

func WriteToFile(filename string, cfg interface{}) error {
	data, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(filename, data, 0666); err != nil {
		return err
	}
	return nil
}
