package websocketHelper

import (
	"fmt"
	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	"github.com/easysoft/zentaoatf/internal/pkg/core/mq"
	"time"
)

var (
	mqTopic  = "WebsocketTopic"
	mqClient *mq.Client
)

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
		msg := mqClient.GetPayLoad(ch).(commDomain.MqMsg)
		fmt.Printf("%s get mq msg '%#v'\n", mqTopic, msg.Content)

		if msg.Content == "exit" {
			mqClient.Unsubscribe(mqTopic, ch)
			break
		} else {
			Broadcast(msg.Namespace, msg.Room, msg.Event, msg.Content)
		}

		time.Sleep(time.Millisecond * 50)
	}
}

func PubMsg(data commDomain.MqMsg) {
	err := mqClient.Publish(mqTopic, data)
	if err != nil {
		fmt.Println("pub mq message failed")
	}
}
