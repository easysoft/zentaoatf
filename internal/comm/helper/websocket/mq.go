package websocketHelper

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/internal/pkg/core/mq"
	"time"
)

var (
	mqTopic  = "WebsocketTopic"
	mqClient *mq.Client
)

type MqMsg struct {
	Namespace string `json:"namespace"`
	Room      string `json:"room"`
	Event     string `json:"event"`
	Content   string `json:"content"`
}

func InitMq() {
	mqClient = mq.NewClient()
	//defer mqClient.Close()
	mqClient.SetConditions(10000)

	go SubMsg()
}

func SubMsg() {
	ch, err := mqClient.Subscribe(mqTopic)
	if err != nil {
		fmt.Printf("sub mq topic %s failed\n", mqTopic)
		return
	}

	for {
		data := mqClient.GetPayLoad(ch).(string)
		fmt.Printf("%s get mq msg '%s'\n", mqTopic, data)

		if data == "exit" {
			mqClient.Unsubscribe(mqTopic, ch)
			break
		} else {
			msg := MqMsg{}
			json.Unmarshal([]byte(data), &msg)
			Broadcast(msg.Namespace, msg.Room, msg.Event, msg.Content)
		}

		time.Sleep(time.Millisecond * 50)
	}
}

func PubMsg(data MqMsg) {
	bytes, _ := json.Marshal(data)

	err := mqClient.Publish(mqTopic, string(bytes))
	if err != nil {
		fmt.Println("pub mq message failed")
	}
}
