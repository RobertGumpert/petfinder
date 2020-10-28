package mailer

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func readFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return bytes, err
}

func headersFromTO(from, to string) string {
	return strings.Join([]string{
		strings.Join([]string{
			"FROM: ", from,
		}, ""),
		strings.Join([]string{
			"TO: ", to,
		}, ""),
	}, "\r\n")
}

func contentTypeHTML() string {
	return strings.Join([]string{
		"Content-type: text/html;charset=utf-8",
		"MIME-Version: 1.0",
	}, "\r\n")
}

func headersSubject(subject string) string {
	return strings.Join([]string{
		"Subject: ", subject,
	}, "")
}

func messageDynamicHTML(to, from, subject string, fileBytes []byte) []byte {
	contentType := contentTypeHTML()
	text := strings.Join([]string{
		headersFromTO(from, to),
		headersSubject(subject),
		contentType,
		string(fileBytes),
	}, "\r\n")
	message := []byte(text)
	return message
}

func messageStaticHTML(to, from, subject, file string) ([]byte, error) {
	fileBytes, err := readFile(file)
	if err != nil {
		return nil, err
	}
	contentType := contentTypeHTML()
	text := strings.Join([]string{
		headersFromTO(from, to),
		headersSubject(subject),
		contentType,
		string(fileBytes),
	}, "\r\n")
	message := []byte(text)
	return message, nil
}

func getTemplate(file string) (*template.Template, error) {
	tmpl, err := template.ParseFiles(file)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}