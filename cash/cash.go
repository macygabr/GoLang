package cash

import (
	"encoding/json"
	"golang/model/task"
	"golang/model/user"
	"io/ioutil"
	"log"

	"github.com/nats-io/stan.go"
)

type Cash struct {
	user user.UserData
}

func NewCash() *Cash {
	return &Cash{}
}

func (c *Cash) Regenerate() stan.Subscription {
	var sub = c.Listen()
	jsonData, err := ioutil.ReadFile("server/download/model.json")
	if err != nil {
		log.Fatal(err)
	}

	if !json.Valid(jsonData) {
		return sub
	}

	err = json.Unmarshal(jsonData, &c.user)
	if err != nil {
		log.Fatal(err)
	}
	return sub
	//Перепиши на нормальное восстановление данных из бд!
}

func (c *Cash) Send(id string) {
	sc, _ := stan.Connect("test-cluster", "client_cash_send", stan.NatsURL("nats://0.0.0.0:4222"))
	defer sc.Close()

	task := new(task.Task)
	task.SetCash(true)
	if c.user.OrderUID == id {
		task.SetUserData(c.user)
	}

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
	})
	return sub
}
