package index

import (
	"github.com/alvan/opsul/api/repos"
)

func init() {
	index("/repos", repos.Index)
}
