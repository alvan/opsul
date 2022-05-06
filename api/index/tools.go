package index

import (
	"github.com/alvan/opsul/api/tools"
)

func init() {
	index("/tools", tools.Index)
}
