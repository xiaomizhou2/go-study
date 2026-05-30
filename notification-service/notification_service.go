package notification_service

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type NotificationSender interface {
	Send(msg string)
}

type EmailSender struct {
	SmtpServer string
}

func (e EmailSender) Send(msg string) error {
	if len(msg) == 0 {
		return fmt.Errorf("邮件内容不能为空")
	}

	time.Sleep(1 * time.Second)

	fmt.Printf("[%s] 邮件已发送: %s", e.SmtpServer, msg)
	return nil
}

type SmsSender struct {
	ApiUrl string
}

func (s SmsSender) Send(msg string) {
	if len(msg) == 0 {
		fmt.Errorf("短信内容不能为空")
	}

	if rand.IntN(100) < 30 {
		fmt.Errorf("短信网关超时: %s", s.ApiUrl)
	}

	fmt.Printf("[%s] 短信已发送: %s", s.ApiUrl, msg)
}
