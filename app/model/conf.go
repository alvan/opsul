package model

import (
	"encoding/json"
	"os"
)

type Conf struct {
	Props map[string]any `json:"props"`
	Https struct {
		Crt string `json:"crt"`
		Key string `json:"key"`
		Use bool   `json:"use"`
	} `json:"https"`
	Files struct {
		Dir string `json:"dir"`
		Use bool   `json:"use"`
	} `json:"files"`
	Webui struct {
		Dir string `json:"dir"`
		Use bool   `json:"use"`
	} `json:"webui"`
	Tmpfs struct {
		Dir string `json:"dir"`
		Pre string `json:"pre"`
	} `json:"tmpfs"`
	Users []*User `json:"users"`
	Repos []*Repo `json:"repos"`
	Tools []*Tool `json:"tools"`
}

// Load from json file
func (this *Conf) Load(file string) (err error) {
	dat, err := os.ReadFile(file)
	if err == nil {
		err = json.Unmarshal(dat, this)
	}
	return
}
