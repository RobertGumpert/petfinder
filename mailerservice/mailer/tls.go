package mailer

import (
	"crypto/tls"
	"net/smtp"
)

// tlsConnect: doc.
//
// Dial подключается к заданному сетевому адресу
// с помощью net.Dial, а затем инициирует рукопожатие TLS,
// возвращая полученное соединение TLS.
//
// NewClient возвращает нового клиента,
// используя существующее соединение и хост
// в качестве имени сервера, которое будет использоваться при аутентификации.
//
func tlsConnect(config *tls.Config, post *post, box *box) (*tls.Conn, *smtp.Client, error) {
	conn, err := tls.Dial("tcp", post.ADDRTls, config)
	if err != nil {
		return nil, nil, err
	}
	client, err := smtp.NewClient(conn, post.Host)
	if err != nil {
		return nil, nil, err
	}
	if err := authSMTPClient(client, box); err != nil {
		//log.Println("ERROR: func tlsConnect : { auth is bad '", err.Error(), "'}")
		if e := conn.Close(); e != nil {
			//log.Println("ERROR: func tlsConnect : { closed connection is bad '", e.Error(), "'}")
			return nil, nil, e
		}
		return nil, nil, err
	}
	return conn, client, nil
}

// authSMTPClient: doc.
//
// useCaseAuthorization аутентифицирует клиента, используя предоставленный механизм аутентификации.
// Неудачная аутентификация закрывает соединение.
// Только серверы, регламинтирующие расширение AUTH, поддерживают эту функцию.
//
func authSMTPClient(client *smtp.Client, box *box) error {
	if err := client.Auth(box.Auth); err != nil {
		return err
	}
	return nil
}

// mailrcpt : doc
//
// Посылает на сервер команду на запуск утилиты mail.
//
// * mail - начинает почтовую транзакцию, которая завершается
// 			передачей данных в один или несколько почтовых ящиков (mail).
//
// * документация - Mail отправляет на сервер команду MAIL,
// 					используя указанный адрес электронной почты.
// 					Если сервер поддерживает расширение 8BITMIME,
// 					Mail добавляет параметр BODY = 8BITMIME.
// 					Это инициирует почтовую транзакцию,
// 					за которой следует один или несколько вызовов Rcpt.
//
// Посылает на сервер команду на запуск утилиты rcpt.
//
// * rcpt - Идентифицирует получателя почтового сообщения.
//
// * документация - Rcpt отправляет на сервер команду RCPT,
// 					используя предоставленный адрес электронной почты.
// 					Вызову Rcpt должен предшествовать вызов Mail,
// 					а за ним может следовать вызов данных или другой вызов Rcpt.
func mailrcpt(client *smtp.Client, sender string, receivers []string) error {
	if err := client.Mail(sender); err != nil {
		return err
	}
	for _, receiver := range receivers {
		if err := client.Rcpt(receiver); err != nil {
			//log.Println("ERROR: func mailrcpt : { client.Rcpt(to.Address) '", err.Error(), "'}")
		}
	}
	return nil
}

// smtpClientLetterBuffer : doc
//
// Посылает на сервер команду на запуск утилиты data.
//
// * data - для отправки текста сообщения. Состоит из заголовка сообщения и тела сообщения,
// 			разделённых пустой строкой.
// 			DATA, по сути, является группой команд,
// 			а сервер отвечает дважды: первый раз на саму команду DATA,
// 			для уведомления о готовности принять текст;
// 			и второй раз после конца последовательности данных,
// 			чтобы принять или отклонить всё письмо..
//
// * документация - Data отправляет на сервер команду DATA
// 					и возвращает средство записи, которое можно использовать
// 					для записи заголовков и тела сообщения.
// 					Вызывающий должен закрыть средство записи перед вызовом
// 					каких-либо других методов в client.
// 					Вызову Data должны предшествовать один или несколько вызовов Rcpt.
func smtpClientLetterBuffer(client *smtp.Client, msg []byte) error {
	writeCloser, err := client.Data()
	if err != nil {
		return err
	}
	_, err = writeCloser.Write(msg)
	if err != nil {
		return err
	}
	err = writeCloser.Close()
	if err != nil {
		return err
	}
	return nil
}

func tlsSendLetter(client *smtp.Client, sender string, receivers []string, msg []byte) error {
	if err := mailrcpt(client, sender, receivers); err != nil {
		return err
	}
	if err := smtpClientLetterBuffer(client, msg); err != nil {
		return err
	}
	if err := client.Quit(); err != nil {
		return err
	}
	return nil
}
