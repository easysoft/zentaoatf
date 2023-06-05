package main

import (
	"fmt"
	"github.com/easysoft/zentaoatf/pkg/core/mq"
	"time"
)

var (
	isExist = false
)

func main() {
	m := mq.NewClient()
	defer m.Close()
	m.SetConditions(10)

	go Sub(m, "exec")

	time.Sleep(time.Second * 3)
	ManyPub(m)
}

func ManyPub(c *mq.Client) {
	for i := 0; i < 100; i++ {
		payload := fmt.Sprintf("msg %02d", i)

		if i == 3 {
			payload = "exit"
		}

		err := c.Publish("exec", payload)
		if err != nil {
			fmt.Println("pub message failed")
			break
		}
		time.Sleep(time.Second * 1)

		if isExist == true {
			break
		}
	}

	c.Close()
}

func Sub(c *mq.Client, top string) {
	ch, err := c.Subscribe(top)
	if err != nil {
		fmt.Printf("sub top:%s failed\n", top)
	}

	for {
		val := c.GetPayLoad(ch)
		fmt.Printf("%s get message is %s\n", top, val)

		if val == "exit" {
			c.Unsubscribe("exec", ch)
			isExist = true
			break
		}
	}
}
