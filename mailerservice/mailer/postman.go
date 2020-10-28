package mailer

import (
	"bytes"
	"crypto/tls"
	"errors"
	"html/template"
	"net/smtp"
	"os"
	"regexp"
	"strings"
)

type letter struct {
	Template *template.Template
	Subject  string
}

type Postman struct {
	Letters map[string]*letter
	Post    *post
	Boxes   map[string]*box
}

type Proto int

const (
	TCP Proto = 0
	TLS Proto = 1
)

type AddNewBox func(host string) (*box, string, error)

func AddBox(key, username, password, identity string) AddNewBox {
	return func(host string) (*box, string, error) {
		if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
			return nil, key, errors.New("User box params must be not equal empty string. ")
		}
		if !EmailValid(username) {
			return nil, key, errors.New("Username not valid. ")
		}
		return &box{
			Auth:     smtp.PlainAuth(identity, username, password, host),
			Username: username,
			Password: password,
			Identity: identity,
		}, key, nil
	}
}

func NewPostman(lettersMap map[string]struct{ FilePath, Subject string }, hostName, tlsPort, tcpPort string, tlsConfig *tls.Config, boxes ...AddNewBox) (*Postman, error) {
	if strings.TrimSpace(hostName) == "" || strings.TrimSpace(tlsPort) == "" || strings.TrimSpace(tcpPort) == "" {
		return nil, errors.New("Connection params must be not equal empty string. ")
	}
	//
	postman := &Postman{
		Letters: make(map[string]*letter, 0),
		Post:    new(post),
		Boxes:   make(map[string]*box, 0),
	}
	//
	post := new(post)
	post.Host = hostName
	post.ADDRTls = strings.Join([]string{hostName, ":", tlsPort}, "")
	post.ADDRTcp = strings.Join([]string{hostName, ":", tcpPort}, "")
	if tlsConfig == nil {
		tlsConfig = &tls.Config{
			ServerName: post.Host,
		}
	}
	post.TLSConfig = tlsConfig
	postman.Post = post
	//
	if len(lettersMap) == 0 {
		return nil, errors.New("Map of template files must be not nil or length not equal 0. ")
	}
	templates := readTemplatesFiles(lettersMap)
	if len(templates) == 0 {
		return nil, errors.New("Length map of template equal 0 after read templates files. ")
	}
	postman.Letters = templates
	//
	if len(boxes) == 0 {
		return nil, errors.New("Array of new boxes must be not nil or length not equal 0. ")
	}
	for _, getNewBox := range boxes {
		box, key, err := getNewBox(hostName)
		if err != nil {
			continue
		}
		postman.Boxes[key] = box
	}
	if len(postman.Boxes) == 0 {
		return nil, errors.New("Length map of boxes equal 0 after call AddNewBox. ")
	}
	//
	return postman, nil
}

func readTemplatesFiles(lettersMap map[string]struct{ FilePath, Subject string }) map[string]*letter {
	templates := make(map[string]*letter, 0)
	for key, selected := range lettersMap {
		if !FileExists(selected.FilePath) {
			continue
		}
		tmpl, err := template.ParseFiles(selected.FilePath)
		if err != nil {
			continue
		}
		templates[key] = &letter{
			Template: tmpl,
			Subject:  selected.Subject,
		}
	}
	return templates
}

func FilledTemplate(tmpl *template.Template, data interface{}) ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := tmpl.Execute(buffer, data)
	return buffer.Bytes(), err
}

func (postman *Postman) TestTLSConnection(keyOfBox string) error {
	var (
		box *box
	)
	if selected, exist := postman.Boxes[keyOfBox]; exist {
		box = selected
	} else {
		return errors.New("Username is not exist. ")
	}
	conn, client, err := tlsConnect(postman.Post.TLSConfig, postman.Post, box)
	if err != nil {
		return err
	}
	if err := conn.Close(); err != nil {
		return err
	}
	if err := client.Quit(); err != nil {
		return err
	}
	return nil
}

func (postman *Postman) SendLetter(keyOfBox, keyOfTemplate, receiver string, data interface{}, proto Proto) error {
	var (
		box *box
		let *letter
	)
	if selected, exist := postman.Boxes[keyOfBox]; exist {
		box = selected
	} else {
		return errors.New("Username is not exist. ")
	}
	if selected, exist := postman.Letters[keyOfTemplate]; exist {
		let = selected
	} else {
		return errors.New("Template is not exist. ")
	}
	tmplBytes, err := FilledTemplate(let.Template, data)
	if err != nil {
		return err
	}
	msg := messageDynamicHTML(receiver, box.Username, let.Subject, tmplBytes)
	switch proto {
	case TLS:
		err = postman.Post.SendLetterTLS(box, []string{receiver}, msg)
		break
	case TCP:
		err = postman.Post.SendLetterTCP(box, []string{receiver}, msg)
		break
	}
	if err != nil {
		return err
	}
	return nil
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

var telephoneRegex = regexp.MustCompile("[+]?[]?[0-9]+-?[0-9]{3}-[0-9]{3}-?[0-9]{4}")

func EmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func TelephoneValid(e string) bool {
	return telephoneRegex.MatchString(e)
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}