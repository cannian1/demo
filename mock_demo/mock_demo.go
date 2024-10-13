//go:generate mockgen -source=./mock_demo.go -destination=./mock_mock_demo.go -package=mock_demo Mail
package mock_demo

import "time"

type Mail interface {
	sendMail(subject, sender, dst, body string) error
}

type MC struct {
	m Mail
}

func NewMC(m Mail) *MC {
	return &MC{m: m}
}

func sign() string {
	return time.Now().Format("2006年01月02日 15:04:05")
}

var getSign = sign

func (c *MC) WriteAndSend(subject, sender, dst, body string) error {
	s := getSign()
	body = body + s
	err := c.m.sendMail(subject, sender, dst, body)
	if err != nil {
		return err
	}
	return nil
}
