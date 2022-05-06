package utils

import (
	"os"
)

type Temp struct {
	Path string `json:"path"`
	Pref string `json:"pref"`
}

func (this *Temp) Init(path string, pref string) (err error) {
	path, err = os.MkdirTemp(path, pref)
	if err == nil {
		this.Path = path
		this.Pref = pref
	}

	return
}

func (this *Temp) File() (*os.File, error) {
	return os.CreateTemp(this.Path, "")
}

func (this *Temp) Done() {
	if this.Path != "" {
		os.RemoveAll(this.Path)
	}
}
