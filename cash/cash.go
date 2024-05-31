package cash

import (
	"encoding/json"
	"golang/model/task"
	"golang/model/user"
	"log"

	"github.com/nats-io/stan.go"
)

type Cash struct {
	// user  user.UserData
	users map[string]user.UserData
}

func NewCash() *Cash {
	return &Cash{}
}

func (c *Cash) Regenerate() stan.Subscription {
	c.users = make(map[string]user.UserData)
	var sub = c.Listen()
	sc, _ := stan.Connect("test-cluster", "client_cash_send", stan.NatsURL("nats://0.0.0.0:4222"))
	defer sc.Close()

	task := new(task.Task)
	task.SetCash(true)

	message, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}
	sc.Publish("database", message)

	return sub
}

func (c *Cash) Send(id string) {
	sc, _ := stan.Connect("test-cluster", "client_cash_send", stan.NatsURL("nats://0.0.0.0:4222"))
	defer sc.Close()

	task := new(task.Task)
	task.SetCash(true)
	task.SetUserData(c.users[id])

	message, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}
	sc.Publish("server", message)
}

func (c *Cash) Listen() stan.Subscription {
	sc, _ := stan.Connect("test-cluster", "client_cash_listen", stan.NatsURL("nats://0.0.0.0:4222"))
	sub, _ := sc.Subscribe("cash", func(msg *stan.Msg) {
		var task task.Task
		err := json.Unmarshal(msg.Data, &task)
		if err != nil {
			log.Fatal(err)
		}
		if task.Cash {
			c.Send(task.OrderID)
		}
		if task.UpdateDB {
			c.users[task.User.OrderUID] = task.User
			// log.Println(c.users)
		}
	})
	return sub
}
