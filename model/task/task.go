package task

import (
	"golang/model/user"
)

type Task struct {
	UpdateDB bool
	Cash     bool
	User     user.UserData
}

func (t *Task) SetUpdateDB(status bool) {
	t.UpdateDB = status
}

func (t *Task) SetCash(status bool) {
	t.Cash = status
}

func (t *Task) SetUserData(user user.UserData) {
	t.User = user
}
