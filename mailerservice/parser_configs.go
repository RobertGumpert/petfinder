package main

import (
	"github.com/spf13/viper"
	"log"
	"mailerservice/mailer"
	"os"
	"path"
	"runtime"
)

func readConfigs(files ...string) map[string]*viper.Viper {
	_, file, _, _ := runtime.Caller(0)
	root := path.Dir(file)
	configs := make(map[string]*viper.Viper)
	var read = func(name string) *viper.Viper {
		vpr := viper.New()
		vpr.SetConfigFile(root + "/" + name + ".yaml")
		if err := vpr.ReadInConfig(); err != nil {
			log.Fatal(err)
		}
		return vpr
	}
	for _, file := range files {
		configs[file] = read(file)
	}
	postman = parse(root, configs)
	return configs
}

func parse(root string, config map[string]*viper.Viper) *mailer.Postman {
	post := config["mail"].GetStringMapString("post")
	letters := make([]struct{ Key, Subject, Template string }, 0)
	lettersMap := make(map[string]struct{ FilePath, Subject string }, 0)
	err := config["mail"].UnmarshalKey("letters", &letters)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	for _, letter := range letters {
		lettersMap[letter.Key] = struct {
			FilePath, Subject string
		}{
			FilePath: root + letter.Template,
			Subject:  letter.Subject,
		}
	}
	boxes := make([]struct{ Key, Username, Password string }, 0)
	boxesArr := make([]mailer.AddNewBox, 0)
	err = config["mail"].UnmarshalKey("boxes", &boxes)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	for _, box := range boxes {
		boxesArr = append(boxesArr,
			mailer.AddBox(
				box.Key,
				box.Username,
				box.Password,
				"",
			),
		)
	}
	postman, err := mailer.NewPostman(
		lettersMap,
		post["host"],
		post["tls"],
		post["tcp"],
		nil,
		boxesArr...,
	)
	return postman
}
