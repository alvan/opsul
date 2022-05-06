package model

type Repo struct {
	Name string  `json:"name"`
	Path string  `json:"path"`
	Pkgs []*Pack `json:"pkgs"`
}

func (this *Repo) FindPackByName(name string) *Pack {
	for _, pack := range this.Pkgs {
		if pack.Name == name {
			return pack
		}
	}

	return nil
}

func (this *Repo) ListNameOfPack() (list []string) {
	for _, pack := range this.Pkgs {
		list = append(list, pack.Name)
	}

	return
}
