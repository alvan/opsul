package index

import (
	"github.com/alvan/opsul/api/tasks"
)

func init() {
	index("/tasks", tasks.Index)
}
