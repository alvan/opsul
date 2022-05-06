package index

import (
	"github.com/alvan/opsul/api/props"
)

func init() {
	index("/props", props.Index)
}
