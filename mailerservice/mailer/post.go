package mailer

import (
	"crypto/tls"
	"errors"
	"net/smtp"
)

type post struct {
	Host         string `json:"host"`
	TLSPort      int    `json:"tls_port"`
	TCPPort      int    `json:"tcp_port"`
	ADDRTcp      string
	ADDRTls      string
	TLSConfig    *tls.Config
}

type box struct {
	Auth     smtp.Auth
	Username string `json:"username"`
	Password string `json:"password"`
	Identity string `json:"identity"`
}

// SendLetterTCP : отправляет от имени агента письмо указанным получателям.
//
// В случае ошибки, пишет ее в канал для ошибок.
//
// * sender - почтовой ящик Агента.
// * receivers - почтовые ящики получателей.
// * message - содержимое письма.
//
func (post *post) SendLetterTCP(sender *box, receivers []string, message []byte) error {
	if err := smtp.SendMail(post.ADDRTcp, sender.Auth, sender.Username, receivers, message); err != nil {
		return errors.New("letter not send [" + err.Error() + "]")
	}
	return nil
}

// SendLetter : отправляет от имени агента письмо указанным получателям.
//
// В случае ошибки, пишет ее в канал для ошибок.
//
// * box - почтовой ящик Агента.
// * receivers - почтовые ящики получателей.
// * message - содержимое письма.
//
func (post *post) SendLetterTLS(sender *box, receivers []string, message []byte) error {
	_, client, err := tlsConnect(post.TLSConfig, post, sender)
	if err != nil {
		return errors.New("smtp.NewClient letter not send [" + err.Error() + "]")
	}
	if err := tlsSendLetter(client, sender.Username, receivers, message); err != nil {
		return errors.New("smtp.NewClient letter not send [" + err.Error() + "]")
	}
	return nil
}
