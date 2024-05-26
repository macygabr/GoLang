package cash

import (
	"encoding/json"
	"golang/model"

	"github.com/nats-io/stan.go"
)

type Cash struct {
	user model.UserData
}

func NewCash() *Cash {
	return &Cash{}
}

func (c *Cash) Regenerate() {
	//read and save user-data from database in cash
}

func (c *Cash) Send() {
	sc, _ := stan.Connect("test-cluster", "client_cash", stan.NatsURL("nats://0.0.0.0:4222"))
	defer sc.Close()
	message, _ := json.Marshal(c.user)
	sc.Publish("UserDataFromCash", message)
}
