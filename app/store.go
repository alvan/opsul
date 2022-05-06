package app

import (
	"sync"

	"github.com/alvan/opsul/app/model"
)

var (
	Store = new(store)
)

type store struct {
	model.Conf
	mutex sync.Mutex

	Tasks []*model.Task
}

func (this *store) AuthBasicUsers() map[string]string {
	users := make(map[string]string)
	for _, item := range this.Users {
		users[item.Name] = item.Pswd
	}

	return users
}

func (this *store) FindUserByName(name string) *model.User {
	for _, item := range this.Users {
		if item.Name == name {
			return item
		}
	}

	return nil
}

func (this *store) FindRepoByName(name string) *model.Repo {
	for _, item := range this.Repos {
		if item.Name == name {
			return item
		}
	}

	return nil
}

func (this *store) FindToolByName(name string) *model.Tool {
	for _, item := range this.Tools {
		if item.Name == name {
			return item
		}
	}

	return nil
}

func (this *store) FindTaskByName(name string) *model.Task {
	for _, item := range this.Tasks {
		if item.Name == name {
			return item
		}
	}

	return nil
}

func (this *store) DropTask(task *model.Task) bool {
	defer this.mutex.Unlock()
	this.mutex.Lock()

	for i, item := range this.Tasks {
		if item.Name == task.Name && item.Id == task.Id {
			this.Tasks = append(this.Tasks[:i], this.Tasks[i+1:]...)
			return true
		}
	}

	return false
}

func (this *store) SaveTask(task *model.Task, otid string) bool {
	defer this.mutex.Unlock()
	this.mutex.Lock()

	for i, item := range this.Tasks {
		if item.Name == task.Name {
			if item.Id != otid {
				return false
			}

			this.Tasks[i] = task
			return true
		}
	}

	if otid == "" {
		this.Tasks = append(this.Tasks, task)
		return true
	}

	return false
}
